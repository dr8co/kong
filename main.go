// kong compiles Monkey source code into bytecode and runs it in a virtual machine.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dr8co/kong/evaluator"
	"github.com/dr8co/kong/lexer"
	"github.com/dr8co/kong/object"
	"github.com/dr8co/kong/parser"
	"github.com/dr8co/kong/repl"
)

const VERSION = "0.1.0"

func main() {
	// Define command-line flags
	fileFlag := flag.String("file", "", "Execute a Monkey script file")
	evalFlag := flag.String("eval", "", "Evaluate a Monkey expression and print the result")
	debugFlag := flag.Bool("debug", false, "Enable debug mode with more verbose output")
	versionFlag := flag.Bool("version", false, "Show version information")

	// Define short flag aliases
	flag.StringVar(fileFlag, "f", "", "Execute a Monkey script file")
	flag.StringVar(evalFlag, "e", "", "Evaluate a Monkey expression and print the result")
	flag.BoolVar(debugFlag, "d", false, "Enable debug mode with more verbose output")
	flag.BoolVar(versionFlag, "v", false, "Show version information")

	// Parse command-line flags
	flag.Parse()

	// Show version information if requested
	if *versionFlag {
		fmt.Printf("Kong Monke Compiler v%s\n", VERSION)
		return
	}

	// Execute a file if specified
	if *fileFlag != "" {
		executeFile(*fileFlag, *debugFlag)
		return
	}

	// Evaluate an expression if specified
	if *evalFlag != "" {
		evaluateExpression(*evalFlag)
		return
	}

	// Get current user
	username := "unknown"
	if usr, err := user.Current(); err == nil {
		username = usr.Username
	}

	fmt.Println("Hello", username+",", "welcome to kong compiler!")
	fmt.Println("Feel free to type in Monkey code. (Ctrl+D or Ctrl+C to exit)")

	// Start the REPL
	repl.Start(os.Stdin, os.Stdout)
}

// executeFile reads and executes a Monkey script file
func executeFile(filename string, debug bool) {
	cleaned := filepath.Clean(filename)
	absolute, err := filepath.Abs(cleaned)
	if err != nil {
		fmt.Printf("Error getting absolute path: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Executing file: %s\n", absolute)

	// Read the file
	//nolint:gosec // We're not reading user input here
	content, err := os.ReadFile(absolute)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Create environment
	env := object.NewEnvironment()

	// Parse and evaluate the file
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)

	// Print the result if in debug mode
	if debug && evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

// evaluateExpression evaluates a single Monkey expression
func evaluateExpression(expr string) {
	// Create environment
	env := object.NewEnvironment()

	// Parse and evaluate the expression
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)

	// Print the result
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

// printParserErrors prints parser errors to stderr
func printParserErrors(errors []string) {
	_, err := fmt.Fprintln(os.Stderr, "Parser errors:")
	if err != nil {
		panic(err)
	}
	for _, msg := range errors {
		_, err := fmt.Fprintln(os.Stderr, "\t"+msg)
		if err != nil {
			panic(err)
		}
	}
}
