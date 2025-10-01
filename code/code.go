// Package code provides bytecode instruction definitions and utilities for the compiler and virtual machine.
//
// This package defines the bytecode instruction set used by the compiler to generate executable code
// and by the virtual machine to execute programs.
//
// It includes opcode definitions, instruction encoding
// and decoding functions, and utilities for working with bytecode instructions.
package code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Instructions is a slice of bytes representing a sequence of instructions.
type Instructions []byte

// Opcode represents a single bytecode instruction used by the compiler and virtual machine.
type Opcode byte

// Bytecode instruction opcodes.
//
// Each opcode represents a specific operation that the virtual machine can execute.
// Instructions may have zero or more operands encoded after the opcode byte.
const (
	// OpConstant pushes a constant from the constant pool onto the stack.
	//
	// Operands: [constant_index:2] - 2-byte index into the constant pool.
	OpConstant Opcode = iota

	// OpAdd pops two values from the stack, adds them, and pushes the result.
	//
	// Stack: [a, b] -> [a + b]
	OpAdd

	// OpPop removes the top value from the stack and discards it.
	//
	// Stack: [value] -> []
	OpPop

	// OpSub pops two values from the stack, subtracts the second from the first, and pushes the result.
	//
	// Stack: [a, b] -> [a - b]
	OpSub

	// OpMul pops two values from the stack, multiplies them, and pushes the result.
	//
	// Stack: [a, b] -> [a * b]
	OpMul

	// OpDiv pops two values from the stack, divides the first by the second, and pushes the result.
	//
	// Stack: [a, b] -> [a / b]
	OpDiv

	// OpTrue pushes the boolean value true onto the stack.
	//
	// Stack: [] -> [true]
	OpTrue

	// OpFalse pushes the boolean value false onto the stack.
	//
	// Stack: [] -> [false]
	OpFalse

	// OpEqual pops two values from the stack, compares them for equality, and pushes the boolean result.
	//
	// Stack: [a, b] -> [a == b]
	OpEqual

	// OpNotEqual pops two values from the stack, compares them for inequality, and pushes the boolean result.
	//
	// Stack: [a, b] -> [a != b]
	OpNotEqual

	// OpGreaterThan pops two values from the stack, compares them, and pushes true if the first is greater.
	//
	// Stack: [a, b] -> [a > b]
	OpGreaterThan

	// OpMinus pops a value from the stack, negates it, and pushes the result.
	//
	// Stack: [value] -> [-value]
	OpMinus

	// OpBang pops a value from the stack, applies logical NOT, and pushes the boolean result.
	//
	// Stack: [value] -> [!value]
	OpBang

	// OpJumpNotTruthy pops a value from the stack and jumps to the specified position if the value is not truthy.
	//
	// Operands: [jump_position:2] - 2-byte absolute instruction position to jump to.
	OpJumpNotTruthy

	// OpJump unconditionally jumps to the specified instruction position.
	//
	// Operands: [jump_position:2] - 2-byte absolute instruction position to jump to.
	OpJump

	// OpNull pushes the null value onto the stack.
	//
	// Stack: [] -> [null]
	OpNull

	// OpGetGlobal retrieves a global variable by index and pushes its value onto the stack.
	//
	// Operands: [global_index:2] - 2-byte index into the global variables store.
	OpGetGlobal

	// OpSetGlobal pops a value from the stack and stores it in the global variable at the specified index.
	//
	// Operands: [global_index:2] - 2-byte index into the global variables store.
	//
	// Stack: [value] -> []
	OpSetGlobal

	// OpArray pops the specified number of elements from the stack and creates an array from them.
	//
	// Operands: [element_count:2] - 2-byte count of elements to pop.
	//
	// Stack: [elem1, elem2, ..., elemN] -> [array]
	OpArray

	// OpHash pops the specified number of key-value pairs from the stack and creates a hash map from them.
	//
	// Operands: [pair_count:2] - 2-byte count of key-value pairs (total stack items = pair_count * 2).
	//
	// Stack: [key1, value1, key2, value2, ..., keyN, valueN] -> [hash]
	OpHash

	// OpIndex pops an index and a collection from the stack, retrieves the element at that index, and pushes it.
	//
	// Stack: [collection, index] -> [collection[index]]
	OpIndex

	// OpCall calls a function with the specified number of arguments.
	//
	// Operands: [num_args:1] - 1-byte count of arguments on the stack.
	//
	// Stack: [func, arg1, arg2, ..., argN] -> [return_value]
	OpCall

	// OpReturnValue pops a value from the stack and returns it from the current function.
	//
	// Stack: [return_value] -> []
	OpReturnValue

	// OpReturn returns from the current function without a return value (implicit null).
	//
	// Stack: [] -> []
	OpReturn

	// OpGetLocal retrieves a local variable by index and pushes its value onto the stack.
	//
	// Operands: [local_index:1] - 1-byte index into the current frame's local variables.
	OpGetLocal

	// OpSetLocal pops a value from the stack and stores it in the local variable at the specified index.
	//
	// Operands: [local_index:1] - 1-byte index into the current frame's local variables.
	//
	// Stack: [value] -> []
	OpSetLocal

	// OpGetBuiltin retrieves a builtin function by index and pushes it onto the stack.
	//
	// Operands: [builtin_index:1] - 1-byte index into the builtin functions table.
	OpGetBuiltin

	// OpClosure creates a closure from a compiled function and captures the specified number of free variables.
	//
	// Operands: [constant_index:2, num_free:1] - 2-byte index to the compiled function in the constant pool,
	// and 1-byte count of free variables to capture from the stack.
	//
	// Stack: [free1, free2, ..., freeN] -> [closure]
	OpClosure

	// OpGetFree retrieves a free variable (captured by a closure) by index and pushes its value onto the stack.
	//
	// Operands: [free_index:1] - 1-byte index into the current closure's free variables.
	OpGetFree

	// OpCurrentClosure pushes the currently executing closure onto the stack (used for recursion).
	//
	// Stack: [] -> [current_closure]
	OpCurrentClosure
)

// Definition represents an instruction definition with its name and operand widths.
type Definition struct {
	// The name of the instruction.
	Name string

	// OperandWidths specifies the number of bytes each operand of an instruction occupies.
	OperandWidths []int
}

// definitions is a map of opcodes to their definitions.
var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpAdd:            {"OpAdd", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpGreaterThan:    {"OpGreaterThan", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpBang:           {"OpBang", []int{}},
	OpJumpNotTruthy:  {"OpJumpNotTruthy", []int{2}},
	OpJump:           {"OpJump", []int{2}},
	OpNull:           {"OpNull", []int{}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}},
	OpHash:           {"OpHash", []int{2}},
	OpIndex:          {"OpIndex", []int{}},
	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
	OpClosure:        {"OpClosure", []int{2, 1}},
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

// Lookup returns the [Definition] for the given [Opcode].
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make creates a byte slice representing an instruction using the provided opcode and operands.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}
	return instruction
}

// String provides a human-readable string representation of the [Instructions], formatted with opcodes and operands.
func (ins Instructions) String() string {
	var out strings.Builder

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			_, _ = fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		_, _ = fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += read + 1
	}

	return out.String()
}

// fmtInstruction formats an instruction with its operands into a human-readable string representation.
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// ReadOperands decodes operands from the specified instructions based
// on the definition and returns them with the total bytes read.
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUint16 decodes the first two bytes of the provided [Instructions] as uint16 in big-endian format.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

// ReadUint8 extracts the first byte from the provided [Instructions] slice and returns it as uint8.
func ReadUint8(ins Instructions) uint8 { return ins[0] }
