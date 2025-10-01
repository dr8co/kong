// Package compiler transforms abstract syntax tree (AST) nodes into bytecode instructions.
//
// This package provides a compiler that traverses an AST produced by the parser and generates
// bytecode instructions that can be executed by a virtual machine.
// The compiler handles expression evaluation, control flow, variable scoping,
// function compilation, and constant management.
//
// # Architecture
//
// The compiler uses a stack-based bytecode generation approach with support for:
//
//   - Multiple compilation scopes for nested functions and closures
//   - Symbol tables for variable resolution (local, global, free, and builtin variables)
//   - Constant pooling for literals and compiled functions
//   - Optimizations such as replacing tail OpPop with OpReturn
//
// # Compilation Process
//
// The compiler works by recursively traversing the AST and emitting bytecode instructions:
//
//  1. Expressions are compiled to push their results onto the stack
//  2. Operators pop operands from the stack and push results
//  3. Variables are resolved through symbol tables and compiled to load/store instructions
//  4. Control flow (if/else) is compiled using conditional and unconditional jumps
//  5. Functions are compiled in separate scopes and stored as constants
//  6. Closures capture free variables from enclosing scopes
//
// # Scoping
//
// The compiler maintains a stack of compilation scopes to support nested functions and closures.
// Each scope has its own instruction sequence and tracks the last two emitted instructions for
// optimization purposes.
// Symbol tables manage variable bindings and support lexical scoping with
// proper closure semantics.
package compiler

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dr8co/kong/ast"
	"github.com/dr8co/kong/code"
	"github.com/dr8co/kong/object"
)

// Compiler is responsible for compiling an AST into bytecode instructions and managing compilation states.
type Compiler struct {
	// Holds the collection of constant values encountered during compilation.
	constants []object.Object

	// symbolTable manages variable bindings and symbol resolution.
	symbolTable *SymbolTable

	// Tracks the current compilation scope and its instruction sequence.
	scopes []CompilationScope

	// scopeIndex tracks the current compilation scope.
	scopeIndex int
}

// Bytecode represents the compiled instructions and constants for a program or function.
type Bytecode struct {

	// Holds the compiled bytecode instructions for a program or function.
	Instructions code.Instructions

	// Contains the constant values used in the bytecode, represented as a slice of objects.
	Constants []object.Object
}

// EmittedInstruction represents a bytecode instruction that has been emitted during compilation.
type EmittedInstruction struct {

	// Opcode represents the specific operation code of the emitted bytecode instruction.
	Opcode code.Opcode

	// Position represents the index or location in the instructions' slice where the bytecode instruction is stored.
	Position int
}

// CompilationScope represents a single layer of compilation containing instructions and metadata about recently emitted instructions.
type CompilationScope struct {

	// Represents the sequence of bytecode instructions for the current compilation scope.
	instructions code.Instructions

	// lastInstruction tracks the most recently emitted bytecode instruction within the current compilation scope.
	lastInstruction EmittedInstruction

	// previousInstruction tracks the second most recently emitted bytecode instruction in the current compilation scope.
	previousInstruction EmittedInstruction
}

// newCompilationScope creates a new compilation scope with an empty instruction sequence.
func newCompilationScope() CompilationScope {
	return CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

// New creates a new compiler instance.
func New() *Compiler {
	symbolTable := NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{newCompilationScope()},
		scopeIndex:  0,
	}
}

// NewWithState creates a new compiler instance with a pre-defined symbol table and constant pool.
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	return &Compiler{
		constants:   constants,
		symbolTable: s,
		scopes:      []CompilationScope{newCompilationScope()},
		scopeIndex:  0,
	}
}

// Compile traverses the given AST node and translates it into bytecode instructions for interpretation.
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)
		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}
		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		symbol := c.symbolTable.Define(node.Name.Value)
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.loadSymbol(symbol)

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := make([]ast.Expression, 0, len(node.Pairs))

		for k := range node.Pairs {
			keys = append(keys, k)
		}

		slices.SortFunc(keys, func(a, b ast.Expression) int {
			return strings.Compare(a.String(), b.String())
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)

	case *ast.FunctionLiteral:
		c.enterScope()
		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name)
		}

		for _, param := range node.Parameters {
			c.symbolTable.Define(param.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}
		fnIndex := c.addConstant(compiledFn)
		c.emit(code.OpClosure, fnIndex, len(freeSymbols))

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))
	}
	return nil
}

// addConstant adds a constant value to the constant pool and returns its index.
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit generates a bytecode instruction with the given opcode and operands,
// adds it to the instruction list, and tracks its position.
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

// setLastInstruction updates the most recent and the previous instruction metadata within the current compilation scope.
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// addInstruction appends the given bytecode instruction to the current scope's instructions and returns its starting position.
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
	return posNewInstruction
}

// Bytecode returns the compiled bytecode containing instructions and constants for a program or function.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// lastInstructionIs checks if the last emitted instruction is of the given opcode.
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// removeLastPop removes the last emitted "pop" instruction from the current compilation scope instructions.
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	newInstruction := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = newInstruction
	c.scopes[c.scopeIndex].lastInstruction = previous
}

// replaceInstruction replaces a sequence of bytecode instructions at the specified position with a new instruction sequence.
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

// changeOperand replaces the operand of an instruction at the specified position with a new provided operand.
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

// currentInstructions retrieves the current compilation scope's bytecode instructions.
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// enterScope initializes a new compilation scope, updates scope tracking, and creates a new enclosed symbol table.
func (c *Compiler) enterScope() {
	scope := newCompilationScope()
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

// leaveScope removes the current compilation scope, updates scope tracking, and restores the outer symbol table.
func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return instructions
}

// replaceLastPopWithReturn modifies the last emitted "pop"
// instruction into a "return value" instruction in the current scope.
func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

// loadSymbol generates bytecode to load the value of a symbol from its associated scope using the symbol's index.
func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(code.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
}
