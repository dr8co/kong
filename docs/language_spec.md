# Monkey Programming Language Specification

This document describes the syntax and semantics of the Monkey programming language.

## 1. Introduction

Monke is a simple, interpreted programming language with a C-like syntax.

It is dynamically typed and supports first-class functions, closures,
and basic data structures like arrays and hash maps.

## 2. Lexical Elements

### 2.1 Comments

Monkey supports single-line comments using the `//` sequence.
Any text from `//` to the end of the line is ignored by the lexer and has no effect on program execution.

Block comments (`/* ... */`) are not supported at this time.

Example:

```monkey
let x = 5; // this is a comment and will be ignored
// let ignored = 10;
puts(x); // prints 5
```

### 2.2 Identifiers

Identifiers start with a letter or underscore and can contain letters, digits, and underscores.

```txt
identifier = letter { letter | digit | "_" } .
letter = "a"..."z" | "A"..."Z" | "_" .
digit = "0"..."9" .
```

### 2.3 Keywords

The following keywords are reserved and cannot be used as identifiers:

```txt
fn    let    true    false    if    else    return
```

### 2.4 Operators and Delimiters

The following characters and character sequences represent operators and delimiters:

```txt
+    -    *    /    =    ==    !=    <    <=    >    >=    !
(    )    {    }    [    ]    ,    ;    :
```

### 2.5 Literals

#### 2.5.1 Integer Literals

Integer literals consist of a sequence of digits.

```txt
integer = digit { digit } .
```

#### 2.5.2 String Literals

String literals are enclosed in double quotes.

```txt
string = '"' { character } '"' .
```

#### 2.5.3 Boolean Literals

Boolean literals are `true` and `false`.

#### 2.5.4 Array Literals

Array literals are enclosed in square brackets and contain a comma-separated list of expressions.

```txt
array = "[" [ expression { "," expression } ] "]" .
```

#### 2.5.5 Hash Literals

Hash literals are enclosed in curly braces and contain a comma-separated list of key-value pairs.

```txt
hash = "{" [ expression ":" expression { "," expression ":" expression } ] "}" .
```

## 3. Types

Monkey has the following built-in types:

- Integer: 64-bit signed integer
- Boolean: true or false
- String: sequence of characters
- Array: ordered collection of values
- Hash: collection of key-value pairs
- Function: first-class function
- Null: represents the absence of a value

## 4. Expressions

### 4.1 Primary Expressions

#### 4.1.1 Identifiers

Identifiers refer to variables, functions, or built-in functions.

#### 4.1.2 Literals

Literals represent fixed values.

#### 4.1.3 Parenthesized Expressions

Expressions can be enclosed in parentheses to control precedence.

```txt
( expression )
```

### 4.2 Function Literals

Function literals define anonymous functions.

```txt
fn ( parameters ) { statements }
```

#### 4.2.1 Closures

Functions in Monkey are first-class values and support lexical scoping.

A closure is a function value that captures variables from its defining environment
and retains access to them even after the outer function returns.

Simple closure example:

```monkey
let newAdder = fn(x) {
 fn(y) { x + y }
};

let addTwo = newAdder(2);
addTwo(3); // => 5
```

Explanation: `newAdder(2)` returns an inner function that captures the outer variable `x`
with value `2`. Calling `addTwo(3)` invokes the inner function and computes `x + y` using the captured `x`.

### 4.3 Call Expressions

Call expressions invoke functions.

```txt
expression ( arguments )
```

### 4.4 Index Expressions

Index expressions access elements of arrays or hashes.

```txt
expression [ expression ]
```

### 4.5 Prefix Expressions

Prefix expressions apply an operator to a single operand.

```txt
operator expression
```

Supported prefix operators:

- `-`: Negation (for integers)
- `!`: Logical NOT (for booleans)

### 4.6 Infix Expressions

Infix expressions apply an operator to two operands.

```txt
expression operator expression
```

Supported infix operators:

- `+`: Addition (for integers and strings)
- `-`: Subtraction (for integers)
- `*`: Multiplication (for integers)
- `/`: Division (for integers)
- `<`: Less than (for integers)
- `>`: Greater than (for integers)
- `<=`: Less than or equal to (for integers)
- `>=`: Greater than or equal to (for integers)
- `==`: Equal to (for all types)
- `!=`: Not equal to (for all types)

### 4.7 If Expressions

If expressions provide conditional evaluation.

```txt
if ( expression ) { statements } [ else { statements } ]
```

## 5. Statements

### 5.1 Expression Statements

Expression statements evaluate an expression and discard the result.

```txt
expression ;
```

### 5.2 Let Statements

Let statements bind a value to an identifier.

```txt
let identifier = expression ;
```

### 5.3 Return Statements

Return statements return a value from a function.

```txt
return expression ;
```

### 5.4 Block Statements

Block statements group multiple statements together.

```txt
{ statements }
```

## 6. Built-in Functions

Monkey provides the following built-in functions:

- `len(arg)`: Returns the length of a string or array
- `first(array)`: Returns the first element of an array
- `last(array)`: Returns the last element of an array
- `rest(array)`: Returns a new array containing all elements except the first
- `push(array, element)`: Returns a new array with the element added to the end
- `puts(args...)`: Prints the arguments to the console

## 7. Evaluation Rules

Monkey uses eager evaluation.

Expressions are evaluated from left to right, with operator precedence determining the order of operations.

## 8. Scoping Rules

Monkey uses lexical scoping.

Variables are visible within the block where they are defined and any nested blocks,
unless shadowed by a variable with the same name in a nested block.

## 9. Error Handling

Monkey does not have explicit error handling mechanisms like try/catch.
Runtime errors result in error objects that terminate execution.
