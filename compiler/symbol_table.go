package compiler

// SymbolScope represents the scope of a symbol within a program, such as global, local, builtin, free, or function.
type SymbolScope string

const (
	// GlobalScope represents the global scope of symbols, typically defining symbols accessible throughout the program.
	GlobalScope SymbolScope = "GLOBAL"

	// LocalScope defines the symbol scope for variables declared within a local function or block.
	LocalScope SymbolScope = "LOCAL"

	// BuiltinScope represents the scope used for predefined or built-in symbols in the program.
	BuiltinScope SymbolScope = "BUILTIN"

	// FreeScope represents the symbol scope for variables that are free,
	// i.e., not locally defined but referenced in a nested context.
	FreeScope SymbolScope = "FREE"

	// FunctionScope represents the scope for function symbols,
	// typically defining variables or symbols within a function context.
	FunctionScope SymbolScope = "FUNCTION"
)

// Symbol represents a named entity within a specific scope and its associated index in the symbol table.
type Symbol struct {
	// The name of the symbol.
	Name string

	// The scope of the symbol.
	Scope SymbolScope

	// The position of the symbol within its respective scope or table.
	Index int
}

// SymbolTable manages variable bindings, symbol definition, and resolution within nested or global scopes.
type SymbolTable struct {
	// Outer represents the parent symbol table, allowing nested scopes to resolve symbols defined in enclosing contexts.
	Outer *SymbolTable

	// store is a map that holds symbol definitions, associating their names with their Symbol metadata.
	store map[string]Symbol

	// numDefinitions tracks the number of symbols defined in the symbol table.
	numDefinitions int

	// FreeSymbols holds a collection of symbols that are referenced but not defined in the current scope,
	// resolved to outer scopes.
	FreeSymbols []Symbol
}

// NewSymbolTable creates a new symbol table with an empty symbol store.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:       make(map[string]Symbol),
		FreeSymbols: []Symbol{},
	}
}

// NewEnclosedSymbolTable creates a new symbol table with its outer field set to the provided enclosing symbol table.
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define adds a new symbol with the given name to the symbol table and assigns it a scope and index.
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve looks up a symbol by name in the current symbol table and, if not found, in enclosing scopes recursively.
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if ok {
			if obj.Scope != GlobalScope && obj.Scope != BuiltinScope {
				free := s.defineFree(obj)
				return free, true
			}
		}
	}
	return obj, ok
}

// DefineBuiltin adds a symbol with a built-in scope to the symbol table using the given index and name.
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

// defineFree adds a free symbol to the FreeSymbols collection and assigns it a FreeScope with a new index.
func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}

	symbol.Scope = FreeScope
	s.store[original.Name] = symbol

	return symbol
}

// DefineFunctionName defines a symbol with function scope and index 0,
// storing it in the symbol table by the given name.
func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = symbol
	return symbol
}
