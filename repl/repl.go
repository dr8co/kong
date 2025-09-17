// Package repl implements the Read-Eval-Print Loop for the Monke programming language.
//
// The REPL provides an interactive interface for users to enter Monke code,
// have it evaluated, and see the results immediately. It uses the Charm libraries
// (Bubbletea, Bubbles, and Lipgloss) to create a modern, user-friendly terminal
// interface with features like syntax highlighting and command history.
//
// Key features:
//   - Interactive command input and execution
//   - Command history tracking
//   - Styled output with different colors for results and errors
//   - Persistent environment across commands
//
// The main entry point is the Start function, which initializes and runs the REPL
// with the given username.
package repl

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dr8co/monke/evaluator"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
	"github.com/dr8co/monke/token"
)

const (
	// Prompt is the default prompt for the REPL
	Prompt = ">> "

	// ContPrompt is the continuation prompt used in multiline input mode within the REPL.
	ContPrompt = ".. "
)

// Options contains configuration options for the REPL
type Options struct {
	NoColor bool // Disable syntax highlighting and colored output
	Debug   bool // Enable debug mode with more verbose output
}

// Start initializes and runs the REPL with the given username and options.
// It creates a new bubbletea program with an initial model and runs it.
// The username is displayed in the welcome message of the REPL.
// If an error occurs while running the program, it is printed to the console.
func Start(username string, options Options) {
	// Start the bubbletea program
	p := tea.NewProgram(initialModel(username, options))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// Styling
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	// Error styles
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	parseErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	runtimeErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF8700")).
				Bold(true)

	errorTipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00"))

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#767676"))

	// Syntax highlighting styles
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	identifierStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	literalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C"))

	operatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	delimiterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9"))

	stringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))
)

// ErrorType represents the type of error that occurred
type ErrorType int

const (

	// NoError indicates that no error occurred, typically used as a default or initial value for error handling.
	NoError ErrorType = iota

	// ParseError indicates an error that occurred during the parsing phase of code evaluation or execution.
	ParseError

	// RuntimeError signifies an error that occurs during the execution of a program, typically at runtime.
	RuntimeError
)

// Custom messages for async evaluation
type evalResultMsg struct {
	output    string
	isError   bool
	errorType ErrorType
	elapsed   time.Duration
}

// The model represents the state of the application
type model struct {
	textInput       textinput.Model
	history         []historyEntry
	env             *object.Environment
	username        string
	evaluating      bool
	currentInput    string
	multilineBuffer string // Buffer for multiline input
	isMultiline     bool   // Flag to indicate if we're in multiline mode
	spinner         spinner.Model
	options         Options
}

// applyStyle applies a lipgloss style to a string, respecting the NoColor option
func (m model) applyStyle(style lipgloss.Style, text string) string {
	if m.options.NoColor {
		return text
	}
	return style.Render(text)
}

// historyEntry represents a single entry in the REPL history
type historyEntry struct {
	input          string
	output         string
	isError        bool
	errorType      ErrorType
	evaluationTime time.Duration // Time taken to evaluate
}

// initialModel creates a new model with default values
func initialModel(username string, options Options) model {
	ti := textinput.New()
	ti.Placeholder = "Enter Monkey code"
	ti.Focus()
	ti.Width = 80
	ti.Prompt = promptStyle.Render(Prompt)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6"))

	return model{
		textInput:       ti,
		history:         []historyEntry{},
		env:             object.NewEnvironment(),
		username:        username,
		evaluating:      false,
		multilineBuffer: "",
		isMultiline:     false,
		spinner:         s,
		options:         options,
	}
}

// Init is the first function that will be called
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// isBalanced checks if brackets, braces, and parentheses are balanced in the input
func isBalanced(input string) bool {
	var stack []rune

	for _, char := range input {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

// evalCmd is a command that evaluates Monkey code asynchronously
func evalCmd(input string, env *object.Environment, debug bool) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		// Process the input
		l := lexer.New(input)
		p := parser.New(l)

		// Debug: Print tokenization time
		var tokenizeTime time.Duration
		if debug {
			tokenizeStart := time.Now()
			program := p.ParseProgram()
			tokenizeTime = time.Since(tokenizeStart)
			var output string
			isError := false
			errorType := NoError
			if len(p.Errors()) != 0 {
				isError = true
				errorType = ParseError
				output = formatParseErrors(p.Errors())

				if debug {
					fmt.Printf("DEBUG: Parse errors: %v\n", p.Errors())
				}
			} else {
				// Debug: Print evaluation time
				evalStart := time.Now()
				evaluated := evaluator.Eval(program, env)
				evalTime := time.Since(evalStart)

				if debug {
					fmt.Printf("DEBUG: Tokenize time: %v\n", tokenizeTime)
					fmt.Printf("DEBUG: Eval time: %v\n", evalTime)
				}

				if evaluated != nil {
					// Check if the result is an error object
					if evaluated.Type() == object.ERROR_OBJ {
						isError = true
						errorType = RuntimeError
						output = formatRuntimeError(evaluated.Inspect())

						if debug {
							fmt.Printf("DEBUG: Runtime error: %s\n", evaluated.Inspect())
						}
					} else {
						output = evaluated.Inspect()

						if debug {
							fmt.Printf("DEBUG: Result type: %s\n", evaluated.Type())
						}
					}
				} else {
					output = "nil"
				}
			}

			elapsed := time.Since(start)

			if debug {
				fmt.Printf("DEBUG: Total execution time: %v\n", elapsed)
			}

			return evalResultMsg{
				output:    output,
				isError:   isError,
				errorType: errorType,
				elapsed:   elapsed,
			}
		}
		// Non-debug path (original code)
		program := p.ParseProgram()

		var output string
		isError := false
		errorType := NoError

		if len(p.Errors()) != 0 {
			isError = true
			errorType = ParseError
			output = formatParseErrors(p.Errors())
		} else {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				// Check if the result is an error object
				if evaluated.Type() == object.ERROR_OBJ {
					isError = true
					errorType = RuntimeError
					output = formatRuntimeError(evaluated.Inspect())
				} else {
					output = evaluated.Inspect()
				}
			} else {
				output = "nil"
			}
		}

		elapsed := time.Since(start)

		return evalResultMsg{
			output:    output,
			isError:   isError,
			errorType: errorType,
			elapsed:   elapsed,
		}
	}
}

// formatError formats error messages.
func (m model) formatError(errorStyle *lipgloss.Style, entry *historyEntry, s *strings.Builder) {
	// Split the output to separate the error message from the tips
	parts := strings.Split(entry.output, "\nTips:")
	if len(parts) > 1 {
		if m.options.NoColor {
			s.WriteString(parts[0])
			s.WriteString("\n")
			s.WriteString("Tips:" + parts[1])
		} else {
			s.WriteString(errorStyle.Render(parts[0]))
			s.WriteString("\n")
			s.WriteString(errorTipStyle.Render("Tips:" + parts[1]))
		}
	} else {
		if m.options.NoColor {
			s.WriteString(entry.output)
		} else {
			s.WriteString(errorStyle.Render(entry.output))
		}
	}
}

// Update handles all the updates to our model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.evaluating {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case evalResultMsg:
		// Evaluation completed
		m.evaluating = false

		// Add to history
		m.history = append(m.history, historyEntry{
			input:          m.currentInput,
			output:         msg.output,
			isError:        msg.isError,
			errorType:      msg.errorType,
			evaluationTime: msg.elapsed,
		})

		m.currentInput = ""
		return m, nil

	case tea.KeyMsg:
		// If we're evaluating, ignore key presses except for Ctrl+C
		if m.evaluating && msg.Type != tea.KeyCtrlC {
			return m, m.spinner.Tick
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				// If we're in multiline mode and the user enters an empty line, evaluate the buffer
				if m.isMultiline {
					if m.multilineBuffer == "" {
						m.isMultiline = false
						return m, nil
					}

					// Start evaluation in the background
					m.evaluating = true
					m.currentInput = m.multilineBuffer
					m.textInput.SetValue("")
					m.isMultiline = false

					// Reset the buffer after evaluation
					buffer := m.multilineBuffer
					m.multilineBuffer = ""

					return m, evalCmd(buffer, m.env, m.options.Debug)
				}
				return m, nil
			}

			// If we're in multiline mode, append the input to the buffer
			if m.isMultiline {
				m.multilineBuffer += "\n" + input
				m.textInput.SetValue("")

				// Check if brackets are now balanced
				if isBalanced(m.multilineBuffer) {
					// Start evaluation in the background
					m.evaluating = true
					m.currentInput = m.multilineBuffer
					m.isMultiline = false

					// Reset the buffer after evaluation
					buffer := m.multilineBuffer
					m.multilineBuffer = ""

					return m, evalCmd(buffer, m.env, m.options.Debug)
				}

				return m, nil
			}

			// Check if the input has balanced brackets
			if !isBalanced(input) {
				// Enter multiline mode
				m.isMultiline = true
				m.multilineBuffer = input
				m.textInput.SetValue("")
				return m, nil
			}

			// Start evaluation in the background
			m.evaluating = true
			m.currentInput = input
			m.textInput.SetValue("")

			return m, evalCmd(input, m.env, m.options.Debug)
		}
	}

	// Only update the text input if we're not evaluating
	if !m.evaluating {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	// Ensure the spinner keeps ticking while evaluating
	if m.evaluating {
		return m, m.spinner.Tick
	}

	return m, cmd
}

// View renders the current UI
func (m model) View() string {
	var s strings.Builder

	// Title
	s.WriteString(m.applyStyle(titleStyle, " Monkey Programming Language REPL "))
	s.WriteString("\n")

	// Welcome message
	if m.username != "" {
		s.WriteString(fmt.Sprintf("\nHello %s! Feel free to type in commands\n", m.username))
	}
	s.WriteString("\n")

	// History
	for _, entry := range m.history {
		// Handle multiline input in history
		lines := strings.Split(entry.input, "\n")
		for i, line := range lines {
			if i == 0 {
				s.WriteString(m.applyStyle(promptStyle, Prompt))
			} else {
				s.WriteString(m.applyStyle(promptStyle, ContPrompt))
			}
			s.WriteString(m.highlightCode(line))
			s.WriteString("\n")
		}

		if entry.isError {
			// Use different styles based on the error type
			switch entry.errorType {
			case ParseError:
				m.formatError(&parseErrorStyle, &entry, &s)
			case RuntimeError:
				m.formatError(&runtimeErrorStyle, &entry, &s)
			default:
				if m.options.NoColor {
					s.WriteString(entry.output)
				} else {
					s.WriteString(errorStyle.Render(entry.output))
				}
			}
		} else {
			if m.options.NoColor {
				s.WriteString(entry.output)
			} else {
				s.WriteString(resultStyle.Render(entry.output))
			}
		}

		// Show evaluation time if it took more than 10 ms
		if entry.evaluationTime > 10*time.Millisecond {
			timeStr := fmt.Sprintf(" (%.2fs)", entry.evaluationTime.Seconds())
			if m.options.NoColor {
				s.WriteString(timeStr)
			} else {
				s.WriteString(historyStyle.Render(timeStr))
			}
		}

		s.WriteString("\n\n")
	}

	// Current evaluation
	if m.evaluating {
		if m.options.NoColor {
			s.WriteString(Prompt)
		} else {
			s.WriteString(promptStyle.Render(Prompt))
		}
		s.WriteString(m.highlightCode(m.currentInput))
		s.WriteString("\n")
		s.WriteString(m.spinner.View())
		s.WriteString(" Evaluating...")
		s.WriteString("\n\n")
	}

	// Show multiline buffer if in multiline mode
	if m.isMultiline && !m.evaluating {
		if m.options.NoColor {
			s.WriteString("Current multiline input:\n")
		} else {
			s.WriteString(historyStyle.Render("Current multiline input:\n"))
		}
		// Instead of splitting by lines, highlight the entire buffer for proper indentation
		s.WriteString(m.highlightCode(m.multilineBuffer))
		s.WriteString("\n")
	}

	// Input
	if !m.evaluating {
		// Set the appropriate prompt based on whether we're in multiline mode
		if m.isMultiline {
			if m.options.NoColor {
				m.textInput.Prompt = ContPrompt
			} else {
				m.textInput.Prompt = promptStyle.Render(ContPrompt)
			}
		} else {
			if m.options.NoColor {
				m.textInput.Prompt = Prompt
			} else {
				m.textInput.Prompt = promptStyle.Render(Prompt)
			}
		}
		s.WriteString(m.textInput.View())
		s.WriteString("\n")
	}

	// Help text
	helpText := "\nPress Esc or Ctrl+C/D to exit"
	if m.isMultiline {
		helpText += " | Multiline mode: Enter empty line to evaluate or continue typing"
	} else {
		helpText += " | Multiline input supported for unbalanced brackets"
	}
	if m.options.NoColor {
		s.WriteString(helpText)
	} else {
		s.WriteString(historyStyle.Render(helpText))
	}

	return s.String()
}

// formatParseErrors formats parser errors into a string with improved readability
func formatParseErrors(errors []string) string {
	var s strings.Builder
	s.WriteString("Parser Errors:\n")

	for i, msg := range errors {
		s.WriteString(fmt.Sprintf("  %d. %s\n", i+1, msg))
	}

	s.WriteString("\nTips:\n")
	s.WriteString("  • Check for missing parentheses, braces, or semicolons\n")
	s.WriteString("  • Verify that all expressions are properly terminated\n")
	s.WriteString("  • Ensure variable names are valid identifiers\n")

	return s.String()
}

// formatRuntimeError formats a runtime error into a string with improved readability
func formatRuntimeError(errorMsg string) string {
	var s strings.Builder
	s.WriteString("Runtime Error:\n")
	s.WriteString("  " + errorMsg + "\n")

	s.WriteString("\nTips:\n")

	// Add specific tips based on common error patterns
	//nolint:gocritic
	if strings.Contains(errorMsg, "identifier not found") {
		s.WriteString("  • Check if the variable is defined before use\n")
		s.WriteString("  • Verify the variable name is spelled correctly\n")
		s.WriteString("  • Make sure the variable is in scope\n")
	} else if strings.Contains(errorMsg, "wrong number of arguments") {
		s.WriteString("  • Check the function call has the correct number of arguments\n")
		s.WriteString("  • Verify the function definition matches its usage\n")
	} else if strings.Contains(errorMsg, "type mismatch") {
		s.WriteString("  • Ensure operands are of compatible types\n")
		s.WriteString("  • Check if you need to convert types before operation\n")
	} else if strings.Contains(errorMsg, "index") {
		s.WriteString("  • Verify array indices are within bounds\n")
		s.WriteString("  • Ensure you're indexing an array or hash\n")
	} else {
		s.WriteString("  • Review your code logic\n")
		s.WriteString("  • Check for type mismatches or undefined variables\n")
		s.WriteString("  • Consider breaking complex expressions into simpler steps\n")
	}

	return s.String()
}

// highlightCode applies syntax highlighting and formatting to Monkey code
//
//nolint:gocyclo
func (m model) highlightCode(code string) string {
	l := lexer.New(code)
	var s strings.Builder

	var tokens []token.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	isKeyword := func(t token.Token) bool {
		switch t.Type {
		case token.FUNCTION, token.LET, token.TRUE, token.FALSE, token.IF, token.ELSE, token.RETURN:
			return true
		}
		return false
	}
	isOperator := func(t token.Token) bool {
		switch t.Type {
		case token.ASSIGN, token.PLUS, token.MINUS, token.BANG, token.ASTERISK, token.SLASH,
			token.LT, token.GT, token.EQ, token.NOT_EQ:
			return true
		}
		return false
	}
	// isIdentifier := func(t token.Token) bool {
	//	return t.Type == token.IDENT
	// }
	isOpenParen := func(t token.Token) bool {
		return t.Type == token.LPAREN
	}
	isCloseParen := func(t token.Token) bool {
		return t.Type == token.RPAREN
	}
	isOpenBrace := func(t token.Token) bool {
		return t.Type == token.LBRACE
	}
	isCloseBrace := func(t token.Token) bool {
		return t.Type == token.RBRACE
	}
	isDelimiter := func(t token.Token) bool {
		switch t.Type {
		case token.COMMA, token.COLON, token.SEMICOLON, token.LPAREN, token.RPAREN,
			token.LBRACE, token.RBRACE, token.LBRACKET, token.RBRACKET:
			return true
		}
		return false
	}

	indentLevel := 0
	atLineStart := true
	for i := range len(tokens) - 1 {
		tok := tokens[i]
		if tok.Type == token.EOF {
			continue
		}
		var prev token.Token
		if i > 0 {
			prev = tokens[i-1]
		}
		next := tokens[i+1]

		// Insert indentation at the start of a new line
		if atLineStart {
			// Don't add indentation or newline if this is an 'else' token following a closing brace
			if tok.Type == token.ELSE && i > 0 && tokens[i-1].Type == token.RBRACE {
				// Skip indentation for 'else' after closing brace
				atLineStart = false
			} else {
				for range indentLevel {
					s.WriteString("  ")
				}
				atLineStart = false
			}
		}

		// Formatting rules (same as before)
		if isKeyword(tok) && tok.Type != token.TRUE && tok.Type != token.FALSE {
			switch tok.Type {
			case token.LET, token.FUNCTION, token.RETURN, token.IF, token.ELSE:
				if m.options.NoColor {
					s.WriteString(tok.Literal)
				} else {
					s.WriteString(keywordStyle.Render(tok.Literal))
				}
				if !isDelimiter(next) && !isOpenBrace(next) && !isOpenParen(next) {
					s.WriteString(" ")
				}
				continue
			}
		}
		if isKeyword(prev) && (prev.Type == token.IF || prev.Type == token.ELSE || prev.Type == token.FUNCTION) && isOpenParen(tok) {
			s.WriteString(" ")
		}
		// if isIdentifier(prev) && isOpenParen(tok) {
		// no space
		// }
		if isOpenBrace(tok) && !isOpenParen(prev) && !isOperator(prev) {
			s.WriteString(" ")
		}
		if isOperator(tok) {
			// Check if this is a prefix operator (like ! or - before an expression)
			isPrefixOp := false
			if (tok.Type == token.BANG || tok.Type == token.MINUS) &&
				(i == 0 || isOpenParen(prev) || isOperator(prev) || isDelimiter(prev)) {
				isPrefixOp = true
			}

			if !isPrefixOp && i > 0 && (!isDelimiter(prev) || isCloseParen(prev)) {
				s.WriteString(" ")
			}

			if m.options.NoColor {
				s.WriteString(tok.Literal)
			} else {
				s.WriteString(operatorStyle.Render(tok.Literal))
			}

			// Add space after the operator only if it's not a prefix operator
			if !isPrefixOp && !isDelimiter(next) && !isCloseParen(next) && !isCloseBrace(next) {
				s.WriteString(" ")
			}
			continue
		}

		// Syntax highlighting
		switch tok.Type {
		case token.FUNCTION, token.LET, token.TRUE, token.FALSE, token.IF, token.ELSE, token.RETURN:
			if m.options.NoColor {
				s.WriteString(tok.Literal)
			} else {
				s.WriteString(keywordStyle.Render(tok.Literal))
			}
		case token.IDENT:
			if m.options.NoColor {
				s.WriteString(tok.Literal)
			} else {
				s.WriteString(identifierStyle.Render(tok.Literal))
			}
		case token.INT:
			if m.options.NoColor {
				s.WriteString(tok.Literal)
			} else {
				s.WriteString(literalStyle.Render(tok.Literal))
			}
		case token.STRING:
			if m.options.NoColor {
				s.WriteString("\"" + tok.Literal + "\"")
			} else {
				s.WriteString(stringStyle.Render("\"" + tok.Literal + "\""))
			}
		case token.ASSIGN, token.PLUS, token.MINUS, token.BANG, token.ASTERISK, token.SLASH,
			token.LT, token.GT, token.EQ, token.NOT_EQ:
			if m.options.NoColor {
				s.WriteString(tok.Literal)
			} else {
				s.WriteString(operatorStyle.Render(tok.Literal))
			}
		case token.COMMA, token.COLON, token.SEMICOLON, token.LPAREN, token.RPAREN,
			token.LBRACE, token.RBRACE, token.LBRACKET, token.RBRACKET:
			// For semicolons, we handle them differently if they follow a closing brace
			//nolint:revive
			if tok.Type == token.SEMICOLON && i > 0 && tokens[i-1].Type == token.RBRACE {
				// Already handled by the special case below
			} else {
				if m.options.NoColor {
					s.WriteString(tok.Literal)
				} else {
					s.WriteString(delimiterStyle.Render(tok.Literal))
				}
			}
		default:
			s.WriteString(tok.Literal)
		}

		// Handle newlines and indentation
		//nolint:staticcheck
		if tok.Type == token.SEMICOLON {
			// If a semicolon follows a closing brace, it was already written
			// Print a newline after semicolon if the next is not EOF or ELSE
			if next.Type != token.EOF && next.Type != token.ELSE {
				s.WriteString("\n")
				atLineStart = true
			}
		} else if tok.Type == token.RBRACE {
			// Check if the next token is a semicolon
			//nolint:gocritic
			if next.Type == token.SEMICOLON {
				// Add the semicolon immediately after the closing brace without a space
				if m.options.NoColor {
					s.WriteString(";")
				} else {
					s.WriteString(delimiterStyle.Render(";"))
				}
			} else if next.Type != token.EOF && next.Type != token.ELSE {
				// No semicolon after brace, add a newline
				s.WriteString("\n")
				atLineStart = true
			} else if next.Type == token.ELSE {
				// Add a single space between closing brace and else
				s.WriteString(" ")
				// Ensure the 'else' is not treated as the start of a new line
				atLineStart = false
			}
		}
		if tok.Type == token.LBRACE {
			// Print a newline after opening brace if next is not closing brace or EOF
			if next.Type != token.RBRACE && next.Type != token.EOF {
				s.WriteString("\n")
				atLineStart = true
			}
			indentLevel++
		}
		if tok.Type == token.RBRACE {
			if indentLevel > 0 {
				indentLevel--
			}
		}
		if tok.Type == token.SEMICOLON && next.Type == token.RBRACE {
			// Don't add an extra newline if the next token is a closing brace
			atLineStart = false
		}

		// Special case: If this is a closing brace and the next token is a semicolon,
		// we've already written the semicolon in the RBRACE case, so skip ahead
		if tok.Type == token.RBRACE && next.Type == token.SEMICOLON {
			// Skip the next token (semicolon) as we've already processed it
			//nolint:ineffassign,wastedassign
			i++
		}
	}

	return s.String()
}
