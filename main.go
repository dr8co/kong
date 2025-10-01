// kong compiles Monkey source code into bytecode and runs it in a virtual machine.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dr8co/kong/compiler"
	"github.com/dr8co/kong/lexer"
	"github.com/dr8co/kong/parser"
	"github.com/dr8co/kong/repl"
	"github.com/dr8co/kong/vm"
)

const version = "0.1.0"

// printUsage displays custom usage information
func printUsage() {
	_, _ = fmt.Fprintf(os.Stderr, `Kong Monkey Compiler v%s

USAGE:
    %s [OPTIONS]

DESCRIPTION:
    Kong compiles Monkey source code into bytecode and runs it in a virtual machine.
    Without any flags, it starts an interactive REPL (Read-Eval-Print-Loop).

OPTIONS:
    -f, --file <path>       Execute a Monkey script file
    -e, --eval <code>       Evaluate a Monkey expression and print the result
    -d, --debug             Enable debug mode with more verbose output
    -v, --version           Show version information
    -h, --help              Show this help message

EXAMPLES:
    # Start interactive REPL
    %s

    # Execute a script file
    %s -f script.monkey
    %s --file script.monkey

    # Evaluate an expression
    %s -e "let x = 5; x * 2"
    %s --eval "puts(\"Hello, World!\")"

    # Execute with debug mode
    %s -f script.monkey -d

`, version, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	// Set custom usage function
	flag.Usage = printUsage

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
		fmt.Printf("Kong Monke Compiler v%s\n", version)
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

	// Parse the file
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	// Compile the program
	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		fmt.Printf("Compilation error: %s\n", err)
		os.Exit(1)
	}

	// Run the bytecode in the VM
	machine := vm.New(comp.Bytecode())
	err = machine.Run()
	if err != nil {
		fmt.Printf("VM error: %s\n", err)
		os.Exit(1)
	}

	// Print the result if in debug mode
	if debug {
		stackTop := machine.LastPoppedStackItem()
		if stackTop != nil {
			fmt.Println(stackTop.Inspect())
		}
	}
}

// evaluateExpression evaluates a single Monkey expression
func evaluateExpression(expr string) {
	// Parse the expression
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	// Compile the program
	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		fmt.Printf("Compilation error: %s\n", err)
		os.Exit(1)
	}

	// Run the bytecode in the VM
	machine := vm.New(comp.Bytecode())
	err = machine.Run()
	if err != nil {
		fmt.Printf("VM error: %s\n", err)
		os.Exit(1)
	}

	// Print the result
	stackTop := machine.LastPoppedStackItem()
	if stackTop != nil {
		fmt.Println(stackTop.Inspect())
	}
}

// printParserErrors prints parser errors to stderr
func printParserErrors(errors []string) {
	_, _ = fmt.Fprintln(os.Stderr, "Parser errors:")
	for _, msg := range errors {
		_, _ = fmt.Fprintln(os.Stderr, "\t"+msg)
	}
}
