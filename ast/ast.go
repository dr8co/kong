// Package ast defines the Abstract Syntax Tree (AST) for the Monke programming language.
//
// The AST represents the structure of a Monke program after it has been parsed.
// It consists of nodes that represent different language constructs such as expressions,
// statements, and literals. The AST is used by the evaluator to execute the program.
//
// Key components:
//   - Node: The base interface for all AST nodes
//   - Statement: Interface for nodes that represent statements (e.g., let, return)
//   - Expression: Interface for nodes that represent expressions (e.g., literals, function calls)
//   - Program: The root node of the AST, containing a list of statements
package ast

import (
	"strings"

	"github.com/dr8co/kong/token"
)

// Node is the base interface for all AST nodes.
// Every node in the AST must implement this interface.
type Node interface {
	// TokenLiteral returns the literal value of the token associated with this node.
	TokenLiteral() string

	// String returns a string representation of the node for debugging and testing.
	String() string
}

// Statement is the interface for all statement nodes in the AST.
// Statements are language constructs that perform actions but don't produce values.
// Examples include let statements, return statements, and expression statements.
type Statement interface {
	Node
	// Marker method to identify statement nodes.
	statementNode()
}

// Expression is the interface for all expression nodes in the AST.
// Expressions are language constructs that produce values.
// Examples include literals, identifiers, function calls, and operators.
type Expression interface {
	Node

	// Marker method to identify expression nodes.
	expressionNode()
}

// Program is the root node of the AST.
// It represents a complete Monke program and contains a list of statements.
type Program struct {
	// The list of statements in the program.
	Statements []Statement
}

// TokenLiteral returns the literal value of the first token in the program.
// If the program has no statements, it returns an empty string.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a string representation of the program.
// It concatenates the string representations of all statements in the program.
func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// An Identifier represents a name in the program, such as a variable or function name.
type Identifier struct {
	// The token containing the identifier.
	Token token.Token

	// The value (name) of the identifier.
	Value string
}

func (id *Identifier) expressionNode() {}

// TokenLiteral returns the literal value of the identifier token.
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

// String returns the value (name) of the identifier.
func (id *Identifier) String() string { return id.Value }

// LetStatement represents a variable binding statement (e.g., "let x = 5;").
type LetStatement struct {
	// The 'let' token.
	Token token.Token

	// The identifier being bound.
	Name *Identifier

	// The expression that produces the value to bind.
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the literal value of the 'let' token.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns a string representation of the let statement.
// Format: "let <identifier> = <expression>;"
func (ls *LetStatement) String() string {
	var out strings.Builder

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement (e.g., "return 5;").
type ReturnStatement struct {
	// The 'return' token.
	Token token.Token

	// The expression that produces the return value.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal value of the 'return' token.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a string representation of the return statement.
// Format: "return <expression>;"
func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents a statement consisting of a single expression.
// For example, function calls can be used as statements.
type ExpressionStatement struct {
	// The first token of the expression.
	Token token.Token

	// The expression itself.
	Expression Expression
}

func (exp *ExpressionStatement) statementNode() {}

// TokenLiteral returns the literal value of the token associated with this statement.
func (exp *ExpressionStatement) TokenLiteral() string { return exp.Token.Literal }

// String returns a string representation of the expression statement.
// It delegates to the String method of the underlying expression.
func (exp *ExpressionStatement) String() string {
	if exp.Expression != nil {
		return exp.Expression.String()
	}
	return ""
}

// IntegerLiteral represents an integer literal expression in the AST.
// For example, the literal "5" in the expression "x + 5".
type IntegerLiteral struct {
	// The token containing the integer literal.
	Token token.Token

	// The actual integer value.
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this integer.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a string representation of the integer literal.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression represents a prefix operator expression in the AST.
// For example, "-5" or "!true" where "-" and "!" are prefix operators.
type PrefixExpression struct {
	// The prefix operator token (e.g., "!").
	Token token.Token

	// The operator (e.g., "!").
	Operator string

	// The expression to the right of the operator.
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this expression.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a string representation of the prefix expression.
// Format: "(<operator><expression>)"
func (pe *PrefixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an infix operator expression in the AST.
// For example, "5 + 5" or "x == y" where "+" and "==" are infix operators.
type InfixExpression struct {
	// The operator token (e.g., "+")
	Token token.Token

	// The expression to the left of the operator.
	Left Expression

	// The operator (e.g., "+")
	Operator string

	// The expression to the right of the operator.
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this expression.
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a string representation of the infix expression.
// Format: "(<left-expression> <operator> <right-expression>)"
func (ie *InfixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean represents a boolean literal expression in the AST.
// For example, "true" or "false".
type Boolean struct {
	// The token containing the boolean literal.
	Token token.Token

	// The actual boolean value.
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this boolean.
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns a string representation of the boolean literal.
func (b *Boolean) String() string { return b.Token.Literal }

// IfExpression represents an if-else expression in the AST.
// For example, "if (x > y) { x } else { y }".
type IfExpression struct {
	// The 'if' token.
	Token token.Token

	// The condition expression.
	Condition Expression

	// The block to execute if the condition is true.
	Consequence *BlockStatement

	// The block to execute if the condition is false (optional)
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this expression.
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a string representation of the `if expression`.
// Format: "if <condition> <consequence> else <alternative>"
func (ie *IfExpression) String() string {
	var out strings.Builder

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// BlockStatement represents a block of statements enclosed in braces.
// For example, "{ statement1; statement2; }".
type BlockStatement struct {
	// The '{' token.
	Token token.Token

	// The statements within the block.
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral returns the literal value of the token associated with this block.
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns a string representation of the block statement.
// It concatenates the string representations of all statements in the block.
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// FunctionLiteral represents a function definition in the AST.
// For example, "fn(x, y) { return x + y; }".
type FunctionLiteral struct {
	// The 'fn' token.
	Token token.Token

	// The function parameters.
	Parameters []*Identifier

	// The function body.
	Body *BlockStatement

	// The name of the function (optional).
	Name string
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this function.
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a string representation of the function literal.
// Format: "fn <namee>(<parameters>) <body>"
func (fl *FunctionLiteral) String() string {
	var out strings.Builder

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString("<" + fl.Name + ">")
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents a function call in the AST.
// For example, "add(1, 2)" or "fn(x, y){ x + y }(1, 2)".
type CallExpression struct {
	// The '(' token.
	Token token.Token

	// The function being called (can be an identifier or function literal).
	Function Expression

	// The arguments passed to the function.
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this call.
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a string representation of the function call.
// Format: "<function>(<arguments>)"
func (ce *CallExpression) String() string {
	var out strings.Builder
	args := make([]string, 0, len(ce.Arguments))

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteral represents a string literal expression in the AST.
// For example, "hello world".
type StringLiteral struct {
	// The token containing the string literal
	Token token.Token

	// The actual string value
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this string.
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String returns a string representation of the string literal.
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// ArrayLiteral represents an array literal expression in the AST.
// For example, "[1, 2 * 2, 3 + 3]".
type ArrayLiteral struct {
	// The '[' token.
	Token token.Token

	// The elements of the array.
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this array.
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String returns a string representation of the array literal.
// Format: "[<elements>]"
func (al *ArrayLiteral) String() string {
	var out strings.Builder

	elems := make([]string, 0, len(al.Elements))
	for _, el := range al.Elements {
		elems = append(elems, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression represents an index expression in the AST.
// For example, "myArray[1]" or "myHash["key"]".
type IndexExpression struct {
	// The '[' token.
	Token token.Token

	// The expression being indexed (array or hash).
	Left Expression

	// The index expression.
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this expression.
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a string representation of the index expression.
// Format: "(<left-expression>[<index-expression>])"
func (ie *IndexExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// HashLiteral represents a hash literal expression in the AST.
// For example, "{key1: value1, key2: value2}".
type HashLiteral struct {
	// The '{' token.
	Token token.Token

	// The key-value pairs in the hash.
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token associated with this hash.
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

// String returns a string representation of the hash literal.
// Format: "{<key1>:<value1>, <key2>:<value2>, ...}"
func (hl *HashLiteral) String() string {
	var out strings.Builder

	pairs := make([]string, 0, len(hl.Pairs))
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
