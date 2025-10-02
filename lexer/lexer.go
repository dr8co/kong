// Package lexer implements the lexical analyzer for the Monke programming language.
//
// The lexer is responsible for breaking down the source code into tokens,
// which are the smallest units of meaning in the language.
// It reads the input character by character and produces a stream of tokens
// that can be processed by the parser.
//
// Key features:
//   - Tokenization of all language elements (keywords, identifiers, literals, operators, etc.)
//   - Handling of whitespace and comments
//   - Error detection for illegal characters
//   - Support for various token types defined in the token package
//   - Optimized for performance with minimal allocations
//
// The main entry point is the New function, which creates a new Lexer instance,
// and the NextToken method, which returns the next token from the input.
package lexer

import (
	"strings"

	"github.com/dr8co/kong/token"
)

// Common tokens that are reused to reduce allocations
var (
	tokenPlus      = token.Token{Type: token.Plus, Literal: "+"}
	tokenMinus     = token.Token{Type: token.Minus, Literal: "-"}
	tokenSlash     = token.Token{Type: token.Slash, Literal: "/"}
	tokenAsterisk  = token.Token{Type: token.Asterisk, Literal: "*"}
	tokenLT        = token.Token{Type: token.Lt, Literal: "<"}
	tokenLTE       = token.Token{Type: token.Lte, Literal: "<="}
	tokenGT        = token.Token{Type: token.Gt, Literal: ">"}
	tokenGTE       = token.Token{Type: token.Gte, Literal: ">="}
	tokenSemicolon = token.Token{Type: token.Semicolon, Literal: ";"}
	tokenColon     = token.Token{Type: token.Colon, Literal: ":"}
	tokenComma     = token.Token{Type: token.Comma, Literal: ","}
	tokenLParen    = token.Token{Type: token.Lparen, Literal: "("}
	tokenRParen    = token.Token{Type: token.Rparen, Literal: ")"}
	tokenLBrace    = token.Token{Type: token.Lbrace, Literal: "{"}
	tokenRBrace    = token.Token{Type: token.Rbrace, Literal: "}"}
	tokenLBracket  = token.Token{Type: token.Lbracket, Literal: "["}
	tokenRBracket  = token.Token{Type: token.Rbracket, Literal: "]"}
	tokenEOF       = token.Token{Type: token.EOF, Literal: ""}
)

// Lexer represents the lexer for the Monke programming language.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	// Pre-allocates a token to reuse for single-character tokens
	singleCharToken token.Token
}

// readChar reads the next character from the input and advances the position.
// It's optimized to minimize checks and operations.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// New creates a new Lexer with the given input string.
// It initializes the lexer, reads the first character, and sets up the token buffer.
func New(input string) *Lexer {
	l := &Lexer{
		input:           input,
		singleCharToken: token.Token{}, // Initialize the token buffer
	}
	l.readChar()
	return l
}

// NextToken reads the next token from the input.
// It skips whitespace, identifies the token type based on the current character,
// and returns a token with the appropriate type and literal value.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			// Use a pre-allocated token for "=="
			l.readChar() // Advance to the next character after '=='
			return token.Token{Type: token.Eq, Literal: string(ch) + string('=')}
		}
		l.readChar() // Advance to the next character after '='
		return token.Token{Type: token.Assign, Literal: "="}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			// Use a pre-allocated token for "!="
			l.readChar() // Advance to the next character after '!='
			return token.Token{Type: token.NotEq, Literal: string(ch) + string('=')}
		}
		l.readChar() // Advance to the next character after '!'
		return token.Token{Type: token.Bang, Literal: "!"}
	case '+':
		l.readChar() // Advance to the next character after '+'
		return tokenPlus
	case '-':
		l.readChar() // Advance to the next character after '-'
		return tokenMinus
	case '/':
		l.readChar() // Advance to the next character after '/'
		return tokenSlash
	case '*':
		l.readChar() // Advance to the next character after '*'
		return tokenAsterisk
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			// advance past '='
			l.readChar()
			return tokenLTE
		}
		l.readChar() // Advance to the next character after '<'
		return tokenLT
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			// advance past '='
			l.readChar()
			return tokenGTE
		}
		l.readChar() // Advance to the next character after '>'
		return tokenGT
	case ';':
		l.readChar() // Advance to the next character after ';'
		return tokenSemicolon
	case ':':
		l.readChar() // Advance to the next character after ':'
		return tokenColon
	case ',':
		l.readChar() // Advance to the next character after ','
		return tokenComma
	case '(':
		l.readChar() // Advance to the next character after '('
		return tokenLParen
	case ')':
		l.readChar() // Advance to the next character after ')'
		return tokenRParen
	case '{':
		l.readChar() // Advance to the next character after '{'
		return tokenLBrace
	case '}':
		l.readChar() // Advance to the next character after '}'
		return tokenRBrace
	case '[':
		l.readChar() // Advance to the next character after '['
		return tokenLBracket
	case ']':
		l.readChar() // Advance to the next character after ']'
		return tokenRBracket
	case '"':
		// readString returns the unescaped content and a bool indicating whether the
		// string was properly terminated (closed by a matching quote).
		lit, ok := l.readString()
		if !ok {
			// unterminated string literal
			l.singleCharToken.Type = token.Illegal
			l.singleCharToken.Literal = "unterminated string"
			return l.singleCharToken
		}
		tok := token.Token{Type: token.String, Literal: lit}
		l.readChar() // Advance to the next character after the closing quote
		return tok
	case 0:
		return tokenEOF
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return token.Token{
				Type:    token.LookupIdent(literal),
				Literal: literal,
			}
		}
		if isDigit(l.ch) {
			return token.Token{
				Type:    token.Int,
				Literal: l.readNumber(),
			}
		}
		// For illegal characters, reuse the single char token
		l.singleCharToken.Type = token.Illegal
		l.singleCharToken.Literal = string(l.ch)
		l.readChar() // Advance to the next character after the illegal character
		return l.singleCharToken
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readNumber reads a number from the input and returns it as a string.
// It's optimized to avoid unnecessary allocations.
func (l *Lexer) readNumber() string {
	position := l.position
	// Fast-forward through digits
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readIdentifier reads an identifier from the input and returns it as a string.
// It's optimized to avoid unnecessary allocations.
func (l *Lexer) readIdentifier() string {
	position := l.position
	// Fast-forward through letters
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhitespace skips any whitespace characters (and comments) in the input.
// It's optimized to use a single loop.
func (l *Lexer) skipWhitespace() {
	// Fast-forward through whitespace and skip `//` line comments.
	for {
		// skip ordinary whitespace
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
			continue
		}

		// skip // comments until the end of the line or EOF
		if l.ch == '/' && l.peekChar() == '/' {
			// consume both '/' characters
			l.readChar()
			l.readChar()
			// advance until newline or EOF
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			// continue the outer loop to handle any whitespace/newline after the comment
			continue
		}

		break
	}
}

// peekChar returns the next character in the input without advancing the position.
// It's optimized to avoid unnecessary checks.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readString reads a string from the input and returns the unescaped content and
// a boolean indicating whether the string was properly terminated (closed by a quote).
func (l *Lexer) readString() (string, bool) {
	var b strings.Builder

	// advance to the first character inside the quotes
	l.readChar()

	for {
		if l.ch == '"' {
			// properly terminated
			return b.String(), true
		}

		if l.ch == 0 {
			// reached EOF without closing quote
			return b.String(), false
		}

		if l.ch == '\\' {
			// consume the backslash and interpret escape
			l.readChar()
			if l.ch == 0 {
				// backslash at EOF — unterminated
				return b.String(), false
			}
			switch l.ch {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '"':
				b.WriteByte('"')
			case '\\':
				b.WriteByte('\\')
			default:
				// Unknown escape: preserve backslash and the char
				b.WriteByte('\\')
				b.WriteByte(l.ch)
			}
		} else {
			b.WriteByte(l.ch)
		}

		l.readChar()
	}
}
