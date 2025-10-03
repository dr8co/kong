// Package object defines the object system for the Monkey programming language.
//
// This package implements the runtime object system that represents values
// during the execution of a Monkey program.
// It defines various types of objects such as integers, booleans, strings,
// arrays, hashes, functions, and built-ins.
//
// Key components:
//   - [Object] interface: The base interface for all runtime values
//   - Various object types ([Integer], [Boolean], [String], [Array], [Hash], [Function], etc.)
//   - [Environment]: Stores variable bindings during execution
//   - [Hashable] interface: For objects that can be used as hash keys
//   - Optimized hash table implementation with key caching for better performance
//
// The evaluator uses the object system to represent and manipulate values
// during program execution.
package object

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/dr8co/kong/ast"
	"github.com/dr8co/kong/code"
)

//nolint:revive
const (
	IntegerObj          = "INTEGER"
	BooleanObj          = "BOOLEAN"
	StringObj           = "STRING"
	NullObj             = "NULL"
	ReturnValueObj      = "RETURN_VALUE"
	ErrorObj            = "ERROR"
	FunctionObj         = "FUNCTION"
	BuiltinObj          = "BUILTIN"
	ArrayObj            = "ARRAY"
	HashObj             = "HASH"
	CompiledFunctionObj = "COMPILED_FUNCTION_OBJ"
	ClosureObj          = "CLOSURE"
)

// Type represents the type of object.
type Type string

// Object is the interface that wraps the basic operations of all Monkey objects.
// All Monkey objects implement this interface.
type Object interface {
	// Type returns the type of the object as a value of Type.
	Type() Type

	// Inspect returns a string representation of the object.
	Inspect() string
}

// Integer represents a Monkey integer value.
type Integer struct {
	Value int64
}

// Type returns the type of the object.
func (i *Integer) Type() Type { return IntegerObj }

// Inspect returns a string representation of the object.
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

// Boolean represents a Monkey boolean value.
type Boolean struct {
	Value bool
}

// Type returns the type of the object.
func (b *Boolean) Type() Type { return BooleanObj }

// Inspect returns a string representation of the object.
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

// String represents a Monkey string value.
type String struct {
	Value string
	// Cache for the hash key to avoid recalculating it
	hashKey *HashKey
}

// Type returns the type of the object.
func (s *String) Type() Type { return StringObj }

// Inspect returns a string representation of the object.
func (s *String) Inspect() string { return s.Value }

// Null represents a Monkey null value.
type Null struct{}

// Type returns the type of the object.
func (n *Null) Type() Type { return NullObj }

// Inspect returns a string representation of the object.
func (n *Null) Inspect() string { return "null" }

// ReturnValue represents a Monkey return value.
type ReturnValue struct {
	Value Object
}

// Type returns the type of the object.
func (rv *ReturnValue) Type() Type { return ReturnValueObj }

// Inspect returns a string representation of the object.
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error represents a Monkey error.
type Error struct {
	Message string
}

// Type returns the type of the object.
func (e *Error) Type() Type { return ErrorObj }

// Inspect returns a string representation of the object.
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function represents a Monkey function.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement

	// Env is the environment in which the function is defined, used to resolve variables during function execution.
	Env *Environment
}

// Type returns the type of the object.
func (f *Function) Type() Type { return FunctionObj }

// Inspect returns a string representation of the object.
func (f *Function) Inspect() string {
	var out strings.Builder
	params := make([]string, 0, len(f.Parameters))

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// BuiltinFunction represents a Monkey builtin function.
type BuiltinFunction func(args ...Object) Object

// Builtin represents a Monkey builtin.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the type of the object.
func (b *Builtin) Type() Type { return BuiltinObj }

// Inspect returns a string representation of the object.
func (b *Builtin) Inspect() string { return "builtin function" }

// Array represents a Monkey array.
type Array struct {
	Elements []Object
}

// Type returns the type of the object.
func (a *Array) Type() Type { return ArrayObj }

// Inspect returns a string representation of the object.
func (a *Array) Inspect() string {
	var out strings.Builder

	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey represents a hash key.
type HashKey struct {
	Type  Type
	Value uint64
}

// HashKey returns the hash key for the object.
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns the hash key for the object.
func (i *Integer) HashKey() HashKey {
	//nolint:gosec
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns the hash key for the object.
func (s *String) HashKey() HashKey {
	// Return the cached hash key if available
	if s.hashKey != nil {
		return *s.hashKey
	}

	// Calculate the hash key
	h := fnv.New64a()
	_, err := h.Write([]byte(s.Value))
	if err != nil {
		return HashKey{Type: ErrorObj, Value: 0}
	}

	// Create and cache the hash key
	hashKey := HashKey{Type: s.Type(), Value: h.Sum64()}
	s.hashKey = &hashKey
	return hashKey
}

// HashPair represents a hash pair.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a Monkey hash.
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the type of the object.
func (h *Hash) Type() Type { return HashObj }

// Inspect returns a string representation of the object.
func (h *Hash) Inspect() string {
	var out strings.Builder

	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Hashable represents an object that can be used as a hash key.
type Hashable interface {
	HashKey() HashKey
}

// CompiledFunction represents a compiled piece of bytecode with its instructions, local variables, and parameters.
type CompiledFunction struct {
	// Represents the bytecode sequence of a compiled function.
	Instructions code.Instructions

	// NumLocals indicates the number of local variables used within the compiled function.
	NumLocals int

	// NumParameters specifies the number of parameters accepted by the compiled function.
	NumParameters int
}

// Type returns the object type of the compiled function, which is [CompiledFunctionObj].
func (c *CompiledFunction) Type() Type { return CompiledFunctionObj }

// Inspect returns a formatted string representation of the CompiledFunction instance, including its memory address.
func (c *CompiledFunction) Inspect() string { return fmt.Sprintf("CompiledFunction[%p]", c) }

// Closure represents a function and its free variables in a virtual machine's execution context.
type Closure struct {
	// Fn is a reference to the compiled function containing the bytecode and metadata for closure execution.
	Fn *CompiledFunction

	// Free holds the objects representing free variables captured by the closure for use during its execution.
	Free []Object
}

// Type returns the type of the object, specifically [ClosureObj] for instances of Closure.
func (c *Closure) Type() Type { return ClosureObj }

// Inspect returns a string representation of the Closure instance, including its memory address.
func (c *Closure) Inspect() string { return fmt.Sprintf("Closure[%p]", c) }
