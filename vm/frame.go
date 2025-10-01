package vm

import (
	"github.com/dr8co/kong/code"
	"github.com/dr8co/kong/object"
)

// Frame represents an execution frame used to track the state of function calls in the virtual machine.
type Frame struct {
	// cl is a reference to an object.Closure,
	// representing a compiled function and its free variables in the execution frame.
	cl *object.Closure

	// ip is the instruction pointer that tracks the current instruction being executed within the frame.
	ip int

	// basePointer is the index in the VM's stack, marking the beginning of the current frame's execution context.
	basePointer int
}

// NewFrame creates a new execution frame for a given closure and base pointer in the virtual machine's stack.
func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basePointer: basePointer}
}

// Instructions retrieves the bytecode instructions of the compiled function associated with the current frame.
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
