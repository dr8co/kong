package object

import "fmt"

// Builtins is a collection of predefined built-in functions available for use within the language.
var Builtins = []struct {
	// The name of the built-in function.
	Name string

	// The definition (and implementation) of the built-in function.
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}

			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
		},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return nil
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
		},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &Array{Elements: newElements}
				}
				return nil
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
		},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				length := len(arg.Elements)
				if length > 0 {
					return arg.Elements[length-1]
				}
				return nil

			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
		},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				length := len(arg.Elements)
				newElements := make([]Object, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]

				return &Array{Elements: newElements}

			default:
				return newError("argument to `push` not supported, got %s", args[0].Type())

			}
		},
		},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		},
		},
	},
}

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// GetBuiltinByName retrieves a built-in function definition by its name from the predefined [Builtins] collection.
//
// It returns a pointer to the corresponding [Builtin] or nil if the name is not found.
func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
