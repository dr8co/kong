package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dr8co/kong/compiler"
	"github.com/dr8co/kong/lexer"
	"github.com/dr8co/kong/parser"
	"github.com/dr8co/kong/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprint(out, PROMPT)
		if err != nil {
			panic(err)
		}
		scanned := scanner.Scan()
		if !scanned {
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

		comp := compiler.New()
		err = comp.Compile(program)
		if err != nil {
			_, err2 := fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			if err2 != nil {
				panic(err2)
			}

			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			_, err2 := fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			continue
		}

		stackTop := machine.StackTop()

		_, err = io.WriteString(out, stackTop.Inspect()+"\n")
		if err != nil {
			panic(err)
		}
	}
}

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
