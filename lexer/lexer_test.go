package lexer

import (
	"testing"

	"github.com/dr8co/kong/token"
)

// TestNextToken tests the functionality of the NextToken method in the Lexer to ensure all tokens are correctly identified.
func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
    x + y;
};
let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10;
10 != 9;

"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.Lparen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.Rparen, ")"},
		{token.Lbrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.Lparen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.Rparen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
		{token.Gt, ">"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.Lparen, "("},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
		{token.Rparen, ")"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Else, "else"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.String, "foobar"},
		{token.String, "foo bar"},
		{token.Lbracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Rbracket, "]"},
		{token.Semicolon, ";"},
		{token.Lbrace, "{"},
		{token.String, "foo"},
		{token.Colon, ":"},
		{token.String, "bar"},
		{token.Rbrace, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestComments ensures that // style line comments are ignored by the lexer
// whether they appear at end-of-line, on their own line, or directly after code.
func TestComments(t *testing.T) {
	input := `let a = 1; // comment
// full line comment
let b = 2; // another
let c = 3;//no space
let d = 4; /////// multiple slashes
let e = "string with // not a comment";
// comment at EOF`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "a"},
		{token.Assign, "="},
		{token.Int, "1"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "b"},
		{token.Assign, "="},
		{token.Int, "2"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "c"},
		{token.Assign, "="},
		{token.Int, "3"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "d"},
		{token.Assign, "="},
		{token.Int, "4"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "e"},
		{token.Assign, "="},
		{token.String, "string with // not a comment"},
		{token.Semicolon, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestCommentBetweenIdentifiers tests tokenization of input containing inline comments between identifiers.
// Verifies token type and literal values, ensuring inline comments are correctly ignored.
func TestCommentBetweenIdentifiers(t *testing.T) {
	input := "a//inline comment\nb"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Ident, "a"},
		{token.Ident, "b"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestCommentBetweenParenthesis(t *testing.T) {
	input := "(//comment\n    x)"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Lparen, "("},
		{token.Ident, "x"},
		{token.Rparen, ")"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestCommentBetweenArrayElements validates the lexer's ability to handle comments between array elements and return correct tokens.
func TestCommentBetweenArrayElements(t *testing.T) {
	input := "[1,//comment\n2]"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Lbracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Rbracket, "]"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestCommentAfterCommaNoSpace tests the lexer for correct handling of comments immediately after a comma without a space.
func TestCommentAfterCommaNoSpace(t *testing.T) {
	input := "a,//c\nb"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Ident, "a"},
		{token.Comma, ","},
		{token.Ident, "b"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestCommentsInComplexConstructs verifies that the lexer correctly handles complex constructs interspersed with comments.
// This includes proper tokenization of functions, arrays, and comments placed between or after constructs.
func TestCommentsInComplexConstructs(t *testing.T) {
	input := `fn(a, // after first arg
    b) { return [1, // in array
    2, 3]; // after array
}; // after function`

	// Expected token sequence ignoring comments
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Function, "fn"},
		{token.Lparen, "("},
		{token.Ident, "a"},
		{token.Comma, ","},
		{token.Ident, "b"},
		{token.Rparen, ")"},
		{token.Lbrace, "{"},
		{token.Return, "return"},
		{token.Lbracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.Comma, ","},
		{token.Int, "3"},
		{token.Rbracket, "]"},
		{token.Semicolon, ";"},
		{token.Rbrace, "}"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestCommentBeforeSemicolon tests the lexing of tokens, including handling inline comments before semicolons.
// It validates the tokens returned by the lexer against the expected types and literals in a structured input.
func TestCommentBeforeSemicolon(t *testing.T) {
	input := `let x = 1 // inline comment
;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "x"},
		{token.Assign, "="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestDivisionFollowedByComment tests the lexer behavior when encountering a division operator followed by a comment.
// Ensures proper differentiation between tokens and validates the type and literal values of each token.
func TestDivisionFollowedByComment(t *testing.T) {
	input := `5 / // divide then comment`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Int, "5"},
		{token.Slash, "/"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestSingleSlashAtEOF validates that the lexer correctly identifies a single slash token followed by an EOF token.
func TestSingleSlashAtEOF(t *testing.T) {
	input := `/`

	l := New(input)

	tok := l.NextToken()
	if tok.Type != token.Slash || tok.Literal != "/" {
		t.Fatalf("expected single slash token, got type=%q literal=%q", tok.Type, tok.Literal)
	}

	tok = l.NextToken()
	if tok.Type != token.EOF {
		t.Fatalf("expected EOF after single slash, got %q", tok.Type)
	}
}

// TestSpacedSlashes tests token parsing for input containing spaced slashes, ensuring correct token type and literal values.
func TestSpacedSlashes(t *testing.T) {
	input := `/ /`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Slash, "/"},
		{token.Slash, "/"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	input := `"hello\nworld" "tab:\tend" "quote:\"inner\"" "backslash:\\"`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.String, "hello\nworld"},
		{token.String, "tab:\tend"},
		{token.String, "quote:\"inner\""},
		{token.String, "backslash:\\"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	input := `"no end`

	l := New(input)

	tok := l.NextToken()
	if tok.Type != token.Illegal {
		t.Fatalf("expected Illegal token for unterminated string, got %q", tok.Type)
	}
	if tok.Literal != "unterminated string" {
		t.Fatalf("expected literal 'unterminated string', got %q", tok.Literal)
	}
}
