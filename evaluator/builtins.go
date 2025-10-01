package evaluator

import (
	"github.com/dr8co/kong/object"
)

// builtins is a map associating built-in function names with their respective *object.Builtin implementations.
var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"puts":  object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
}
