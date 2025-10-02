# Kong — compiler & VM for the Monkey language

Kong is a small, self-contained compiler and bytecode virtual machine
that implements the Monkey programming language.

Kong compiles Monkey source code to compact bytecode and executes it with a stack-based VM.
A built-in REPL lets you experiment with the language interactively.

It is a faster, more efficient, and more robust
[Monke](https://github.com/dr8co/monke "Monke Interpreter").

## Features

- **Lexer**: Tokenizes input source code.
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens.
- **AST**: Represents the structure of parsed code.
- **Compiler**: Translates the AST to bytecode.
- **Virtual Machine (VM)**: Stack-based VM that executes the bytecode.
- **REPL**: Interactive shell for running Monkey code.
- **Built-in Functions**: Includes basic built-in functions for convenience.
- **First-class Functions**: Supports functions as first-class citizens, including closures.
- **Data Structures**: Supports arrays and hash maps.
- **Error Handling**: Graceful handling of syntax and runtime errors.

## Getting Started

### Prerequisites

- Go 1.25 or newer

### Installation

To install Kong, you can install it via `go install`:

```bash
go install github.com/dr8co/kong@latest
```

This will install the `kong` executable in your `$GOPATH/bin` directory.

You can also clone the repository and build it manually:

```bash
git clone https://github.com/dr8co/kong.git
cd kong
go build
```

**Pre-built binaries** are available for download.

See the [Releases Page](https://github.com/dr8co/kong/releases).

## Project Structure

- `main.go` — CLI & entry point. Starts the REPL by default.
- `token/` — Token definitions.
- `lexer/` — Lexical analyzer (tokenizer).
- `parser/` — Parses Monkey source code into an Abstract Syntax Tree (AST).
- `ast/` — Abstract Syntax Tree definitions.
- `object/` — runtime value representations and environment.
- `compiler/` — Compiler that generates bytecode.
- `code/` — Bytecode instruction definitions and helpers.
- `vm/` — Virtual Machine that executes bytecode.
- `repl/` — the REPL that wires compiler + VM to provide a persistent interactive session.
- `docs/` — design docs, language spec, REPL guide and examples.

## Example Usage

Evaluate a single expression from the command line:

```bash
kong -e 'let x = 5; x + 10;'
```

Run `kong -h` for help and options.

To start the REPL, run `kong` with no arguments:

```bash
kong  # if you installed via go install, ensure $GOPATH/bin is in your PATH
# or
./kong # if you built manually or downloaded the binary
```

### Example REPL Session

```console
$ kong
>> let x = 5;
5
>> x + 10;
15
>> let add = fn(a, b) { a + b; };
Closure[0xc000200000]
>> add(2, 3);
5
>> let arr = [1, 2, 3];
[1, 2, 3]
>> arr[1];
2
>> last(arr);
3
>> puts("Hello, Kong!");
Hello, Kong!
null
```

## Testing

To run the tests:

```bash
go test ./...
```

## Documentation

- Language specification: `docs/language_spec.md` (Monkey syntax & semantics)
- Architecture & design decisions: `docs/architecture.md` (Kong compiler and VM)
- REPL guide & examples: `docs/repl_guide.md` and `docs/examples/`

For more details and examples, check the [documentation](./docs/README.md).

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements or bug fixes.
See [CONTRIBUTIONS.md](./CONTRIBUTIONS.md) for guidelines.

## License

This project is licensed under the MIT License.
See the [LICENSE](./LICENSE) file for details.
