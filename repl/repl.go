// Package repl provides a Read-Eval-Print Loop (REPL) for interactive code execution.
//
// This package implements an interactive shell that allows users to enter code line-by-line,
// see immediate results, and maintain state across multiple inputs.
// The REPL integrates the lexer, parser, compiler, and virtual machine to provide
// a complete execution environment.
//
// # Architecture
//
// The REPL operates in a continuous loop that:
//
//  1. Reads a line of input from the user
//  2. Lexes and parses the input into an abstract syntax tree (AST)
//  3. Compiles the AST into bytecode instructions
//  4. Executes the bytecode in the virtual machine
//  5. Prints the result of the evaluation
//
// # State Management
//
// The REPL maintains persistent state across inputs to support variable declarations
// and function definitions that span multiple interactions:
//
//   - Constants: A growing pool of immutable values compiled from literals
//   - Globals: A fixed-size store for global variables accessible across inputs
//   - Symbol Table: Tracks variable names and their scopes (builtin, global, local)
//
// This allows users to define variables and functions in one input and reference them
// in subsequent inputs, creating a natural interactive programming experience.
//
// # Error Handling
//
// The REPL provides user-friendly error messages for:
//
//   - Parse errors: Syntax errors with context about what went wrong
//   - Compilation errors: Issues during bytecode generation
//   - Runtime errors: Execution failures in the virtual machine
//
// When an error occurs, the REPL displays the error message and continues running,
// allowing users to correct their input and try again without restarting the session.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/dr8co/kong/compiler"
	"github.com/dr8co/kong/lexer"
	"github.com/dr8co/kong/object"
	"github.com/dr8co/kong/parser"
	"github.com/dr8co/kong/vm"
)

// PROMPT is the string used to prompt the user for input.
const PROMPT = ">> "

// Start starts the REPL and runs the interactive loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var constants []object.Object
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		_, err := fmt.Fprint(out, PROMPT)
		if err != nil {
			panic(err)
		}
		scanned := scanner.Scan()
		if !scanned {
			if out == os.Stdout || out == os.Stderr {
				_, _ = fmt.Fprintln(out, "bye!")
			}
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			_, err2 := fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			_, err2 := fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			continue
		}

		lastPopped := machine.LastPoppedStackItem()

		_, err = io.WriteString(out, lastPopped.Inspect()+"\n")
		if err != nil {
			panic(err)
		}
	}
}

// printParseErrors prints a list of parse errors to the given output stream.
func printParseErrors(out io.Writer, errors []string) {
	_, err := io.WriteString(out, "parser errors:\n")
	if err != nil {
		panic(err)
	}

	for _, msg := range errors {
		_, err = io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			panic(err)
		}
	}
}
