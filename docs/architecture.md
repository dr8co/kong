# Kong: Architecture and Design Decisions

This document describes the architecture of the Kong compiler and runtime
and explains the key design decisions made during its development.

## Overall Architecture

Kong uses a compiler + virtual machine (bytecode) execution model:

Source Code (Monkey) -> Lexer -> Parser -> AST -> Compiler -> Bytecode -> VM -> Result

All components (lexer, parser, AST, and object model) are shared across the toolchain.
The REPL and CLI compile inputs to bytecode and execute them on the virtual machine (VM).

The bytecode VM is the primary runtime and provides the performance characteristics required for real workloads.

Each component has a specific responsibility:

1. **Lexer**: Converts source code text into tokens
2. **Parser**: Transforms tokens into an Abstract Syntax Tree (AST)
3. **Compiler**: Compiles the AST to bytecode
4. **VM**: Executes bytecode produced by the compiler

## Component Breakdown

### Lexer (`lexer` package)

The lexer is responsible for breaking down the source code into tokens.
It reads the input character by character and produces a stream of tokens.

**Design Decisions:**

- **Simple Character-by-Character Scanning**: The lexer uses a straightforward approach of reading one character at a time, which makes the code easy to understand and maintain.
- **Look-Ahead**: The lexer uses a one-character look-ahead to handle multi-character tokens like `==` and `!=`.
- **No Regex**: The lexer avoids using regular expressions to remain portable and transparent.
- **Token Types**: Tokens are categorized (keywords, identifiers, literals, operators, etc.) to simplify parsing.

### Parser (`parser` package)

The parser converts the token stream into an AST.
It implements a recursive descent parser with Pratt parsing (precedence climbing) for expressions.

**Design Decisions:**

- **Pratt Parsing**: The parser uses Pratt parsing for expression parsing to handle operator precedence cleanly.
- **Recursive Descent**: For statements and other constructs, recursive descent is used to mirror the language grammar.
- **Error Reporting**: The parser collects multiple errors during a parse to present more helpful diagnostics.
- **Prefix and Infix Functions**: The parser uses maps of prefix and infix parsing functions to handle different types of expressions, making it easy to extend with new expression types.

### AST (`ast` package)

The Abstract Syntax Tree (AST) represents the structure of a program.
Each node in the tree corresponds to a language construct (expression, statement, etc.).

**Design Decisions:**

- **Node Interface**: A unified `Node` interface allows generic handling of AST nodes.
- **Expression and Statement Types**: Clear separation of expressions and statements simplifies evaluation and compilation.
- **String Representation**: Nodes can be rendered to strings for debugging and tests.
- **Immutable Nodes**: AST nodes are designed to be immutable, simplifying the evaluation process.

### Runtime model

Kong compiles source code to bytecode which is executed on a stack-based virtual machine.

The runtime provides:

- Lexical scoping via environments and frames
- First-class functions and closures (functions capture their defining environment)
- Errors represented as runtime objects
- Built-in functions implemented in Go and surfaced to Monkey programs

### Object System (`object` package)

The object system defines the runtime values that can exist in a Monkey program.
It includes integers, booleans, strings, arrays, hashes, functions, and more.

**Design Decisions:**

- **Object Interface**: Runtime values implement a common interface to be treated uniformly.
- **Value Representation**: Each type of value has its own struct with appropriate fields.
- **First-Class Functions**: Functions are treated as first-class values, allowing them to be passed around, returned from other functions, and stored in variables.
- **Hashable Interface**: Types that can be keys in hash maps implement a hashable contract.
- **Environment**: The runtime environment maps identifiers to objects and supports nesting.
- **Closures**: Functions capture their defining environment, enabling closures.
- **Error Handling**: Errors are represented as values that can be passed around, allowing for consistent error handling throughout the evaluation process.
- **Built-in Functions**: Common functions are provided as built-ins, implemented directly in Go rather than in Monkey.

### Compiler (`compiler`, `code` packages)

The compiler translates the AST into bytecode instructions that the VM can execute.
It performs a single pass over the AST, emitting instructions and managing scopes.

**Design Decisions:**

- **Bytecode Format**: Instructions are encoded compactly with operands; constants live in a constants' pool.
- **Scopes and Symbol Tables**: The compiler maintains symbol tables for variable/function resolution, supporting nested scopes.
- **Function Compilation**: Functions are compiled into their own bytecode chunks, allowing for recursion and closures.

### Virtual Machine (`vm` package`)

The VM executes the bytecode produced by the compiler.
It provides a stack-based execution model and manages the runtime environment.

**Design Decisions:**

- **Stack-Based Execution**: The VM uses a stack and frames to implement calls and local state.
- **Frame Management**: Each function call creates a new frame on the stack, allowing for nested calls and proper scoping.
- **Error Handling**: Runtime errors are represented as objects and surfaced in ways that help debugging and testing.

### REPL (`repl` package)

The REPL (Read-Eval-Print Loop) provides an interactive interface for users
to enter Monkey code and see the results immediately.

The REPL compiles inputs to bytecode and executes them on the VM.

**Design Decisions:**

- **Simple Terminal UI**: Clear prompt and reliable behavior is prioritized.
- **Persistent State**: Globals, constants, and symbol tables can persist across inputs.

## Key Design Principles

### Simplicity Over Performance

Kong's implementation emphasizes clarity and modularity, making it a good platform for learning and experimentation.

> [!Note]
> Despite the simplicity, Kong is significantly more performant than
> [the original Monkey interpreter](https://github.com/dr8co/Monke)
> due to the bytecode compilation and VM execution path.

### Modularity

Clear separation of concerns (lexer, parser, AST, compiler, VM, object system)
makes the codebase easier to test and extend.

### Extensibility

The architecture is designed so new language constructs or runtime features can be added with localized changes.

### Error Handling

Errors are represented as runtime objects and surfaced in ways that help debugging and testing.

## Future Directions

Kong currently uses a bytecode compiler+VM as its primary execution model.

Possible directions include:

1. **JIT Compilation**: Add a JIT for hot code paths.
2. **Static Analysis/Type Checking**: Optional static checks to catch issues earlier.
3. **Module System & Stdlib**: Improve code organization and grow a standard library of built-ins.

## Conclusion

Kong combines a clear, testable implementation with a higher-performance bytecode path.

The design balances learnability with practical execution features,
making it a good foundation for language experimentation and extension.
