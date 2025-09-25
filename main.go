package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dr8co/kong/repl"
)

func main() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Println("Hello", usr.Username, "This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")

	repl.Start(os.Stdin, os.Stdout)
}
