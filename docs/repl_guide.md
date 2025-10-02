# Kong REPL User Guide

This guide explains how to use the Kong Read-Eval-Print Loop (REPL) interface,
which compiles your input to bytecode and executes it on the Kong virtual machine.

## Getting Started

To start the Kong REPL, run the `kong` executable:

```bash
kong
```

You'll see a welcome message and a prompt (`>>`) where you can start entering Monkey code.

## Basic Usage

Type expressions or statements at the prompt. Example:

```console
>> let x = 5;
5
```

Multi-line input is supported. The REPL evaluates once the expression or block is complete:

```console
>> if (x > 3) {
  x + 10
} else {
  x - 10
}
15
```

You can define variables and functions, and they persist in the session:

```console
>> let name = "World";
>> let greeting = "Hello, " + name + "!";
>> greeting
"Hello, World!"

>> let add = fn(a, b) { a + b };
>> add(5, 10)
15
```

Arrays and hash maps are supported:

```console
>> let myArray = [1, 2, 3, 4, 5];
>> myArray[2]
3

>> let myHash = {"name": "Monkey", "type": "Language"};
>> myHash["name"]
"Monkey"
```

Built-in functions include `len`, `first`, `last`, `rest`, `push`, and `puts`:

```console
>> len("hello")
5

>> puts("Hello, World!")
Hello, World!
null
```

## Keyboard Shortcuts

- **Enter**: Execute the current input
- **EOF** (**Ctrl+D** on Unix, **Ctrl+Z+Enter** on Windows) or **Ctrl+C**: Exit the REPL

## Tips

- **Persistent Environment**: Variables and functions defined in the REPL persist for the session.
- **Semicolons**: Statements should end with a semicolon (`;`).

## Example Session

```console
>> let x = 10;
10

>> let y = 5;
5

>> x + y
15

>> let max = fn(a, b) { if (a > b) { a } else { b } };
Closure[0xc0000720c0]

>> max(x, y)
10

>> let array = [1, 2, 3, 4, 5];
[1, 2, 3, 4, 5]

>> len(array)
5

>> array[2]
3
```

For more information about the Monkey language, see the
[language specification](./language_spec.md).
