// Package token defines the token types and structures for the Monke programming language.
//
// Tokens are the smallest units of meaning in the language, produced by the lexer
// during the lexical analysis phase.
// Each token represents a specific language element such as a keyword,
// identifier, operator, or delimiter.
//
// Key components:
//   - [Type]: A type representing different categories of tokens
//   - [Token]: A structure containing the type and literal value of a token
//   - Constants for all token types supported by the language
//   - Lookup functions for identifying keywords
//
// This package is used primarily by the lexer to categorize input text and by the
// parser to understand the structure of the program.
package token

// Type represents the type of token.
type Type string

// Token represents a single token in the source code.
type Token struct {
	// Type specifies the category of the token, such as keywords, identifiers, or operators.
	Type Type

	// Literal specifies the exact string value of the token as it appears in the source code.
	Literal string
}

const (
	// Single-character tokens

	// Illegal represents an unknown or invalid token that the lexer does not recognize.
	Illegal = "Illegal"

	// EOF represents the end-of-file marker, signaling that no more tokens are available.
	EOF = "EOF"

	// Identifiers & literals

	// Ident represents an identifier token, such as variable names or function names.
	Ident = "Ident"

	// Int represents an integer literal token.
	Int = "Int"

	// String represents a string literal token.
	String = "String"

	// Operators

	// Assign represents the assignment operator "=".
	Assign = "="

	// Plus represents the addition operator "+".
	Plus = "+"

	// Minus represents the subtraction operator "-".
	Minus = "-"

	// Bang represents the logical NOT operator "!".
	Bang = "!"

	// Asterisk represents the multiplication operator "*".
	Asterisk = "*"

	// Slash represents the division operator "/".
	Slash = "/"

	// Lt represents the less-than comparison operator "<".
	Lt = "<"

	// Lte represents the less-than-or-equal-to comparison operator "<=".
	Lte = "<="

	// Gt represents the greater-than comparison operator ">".
	Gt = ">"

	// Gte represents the greater-than-or-equal-to comparison operator ">=".
	Gte = ">="

	// Eq represents the equality comparison operator "==".
	Eq = "=="

	// NotEq represents the inequality comparison operator "!=".
	NotEq = "!="

	// Delimiters

	// Comma represents the comma delimiter ",".
	Comma = ","

	// Colon represents the colon delimiter ":".
	Colon = ":"

	// Semicolon represents the semicolon delimiter ";".
	Semicolon = ";"

	// Lparen represents the left parenthesis delimiter "(".
	Lparen = "("

	// Rparen represents the right parenthesis delimiter ")".
	Rparen = ")"

	// Lbrace represents the left brace delimiter "{".
	Lbrace = "{"

	// Rbrace represents the right brace delimiter "}".
	Rbrace = "}"

	// Lbracket represents the left bracket delimiter "[".
	Lbracket = "["

	// Rbracket represents the right bracket delimiter "]".
	Rbracket = "]"

	// Keywords

	// Function represents the "fn" keyword for function declarations.
	Function = "Function"

	// Let represents the "let" keyword for variable declarations.
	Let = "Let"

	// True represents the "true" boolean literal keyword.
	True = "True"

	// False represents the "false" boolean literal keyword.
	False = "False"

	// If represents the "if" keyword for conditional expressions.
	If = "If"

	// Else represents the "else" keyword for alternative branches in conditional expressions.
	Else = "Else"

	// Return represents the "return" keyword for returning values from functions.
	Return = "Return"
)

// keywords is a map of reserved keywords to their corresponding token types.
var keywords = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// LookupIdent checks if the given identifier is a keyword.
// If it is, it returns the corresponding token type.
// Otherwise, it returns the [Ident] type.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
