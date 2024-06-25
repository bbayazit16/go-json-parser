package main

import (
	"math"
	"strings"
)

type Lexer struct {
	input   []rune
	tokens  []Token
	line    int
	start   int
	current int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		line:  1,
	}
}

func (lexer *Lexer) Scan() ([]Token, error) {
	for !lexer.isAtEnd() {
		lexer.start = lexer.current

		lexerError := lexer.scanToken()
		if lexerError != nil {
			return nil, lexerError
		}

	}

	return lexer.tokens, nil
}

func (lexer *Lexer) scanToken() *LexerError {
	char := lexer.advance()
	switch char {
	case '{':
		lexer.addToken(LBRACE, lexer.line, lexer.start, "{")
	case '}':
		lexer.addToken(RBRACE, lexer.line, lexer.start, "}")
	case '[':
		lexer.addToken(LBRACKET, lexer.line, lexer.start, "[")
	case ']':
		lexer.addToken(RBRACKET, lexer.line, lexer.start, "]")
	case ',':
		lexer.addToken(COMMA, lexer.line, lexer.start, ",")
	case ':':
		lexer.addToken(COLON, lexer.line, lexer.start, ":")

	case '"':
		lexerError := lexer.string()
		if lexerError != nil {
			return lexerError
		}

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		lexerError := lexer.number()
		if lexerError != nil {
			return lexerError
		}

	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		lexer.line++

	case 't', 'f', 'n':
		lexerError := lexer.keyword()
		if lexerError != nil {
			return lexerError
		}

	default:
		return &LexerError{
			Line:      lexer.line,
			Character: lexer.current,
			Message:   "Unexpected character: " + string(char),
		}
	}

	return nil
}

func (lexer *Lexer) keyword() *LexerError {
	for !lexer.isAtEnd() && (lexer.peek() >= 'a' && lexer.peek() <= 'z') {
		lexer.advance()
	}

	value := string(lexer.input[lexer.start:lexer.current])

	var tokenType TokenType
	switch value {
	case "true":
		tokenType = BOOL
	case "false":
		tokenType = BOOL
	case "null":
		tokenType = NULL
	default:
		return &LexerError{
			Line:      lexer.line,
			Character: lexer.current,
			Message:   "Invalid keyword " + value,
		}
	}

	lexer.addToken(tokenType, lexer.line, lexer.start, value)

	return nil
}

func (lexer *Lexer) number() *LexerError {
	for lexer.peek() >= '0' && lexer.peek() <= '9' {
		lexer.advance()
	}

	if lexer.peek() == '.' {
		// Consume .
		lexer.advance()

		if lexer.peek() < '0' || lexer.peek() > '9' {
			return &LexerError{
				Line:      lexer.line,
				Character: lexer.current,
				Message:   "Invalid number. Expected digit after .",
			}
		}

		for lexer.peek() >= '0' && lexer.peek() <= '9' {
			lexer.advance()
		}
	}

	value := string(lexer.input[lexer.start:lexer.current])
	lexer.addToken(NUMBER, lexer.line, lexer.start, value)

	return nil
}

func (lexer *Lexer) string() *LexerError {
	var builder strings.Builder

	for !lexer.isAtEnd() {
		nextChar := lexer.peek()

		if nextChar == '"' {
			break
		}

		if lexer.isAtEnd() {
			return &LexerError{
				Line:      lexer.line,
				Character: lexer.current,
				Message:   "Unterminated string",
			}
		}

		if nextChar == '\\' {
			// Consume backslash
			lexer.advance()

			nextChar = lexer.peek()
			switch nextChar {
			// All escapable characters
			case 'b', 'f', 'n', 'r', 't':
				builder.WriteRune('\\' + lexer.advance())
			case '\\', '"', '/':
				builder.WriteRune(lexer.advance())
			case 'u':
				unicode, err := lexer.unicode()
				if err != nil {
					return err
				}
				builder.WriteRune(unicode)
			default:
				return &LexerError{
					Line:      lexer.line,
					Character: lexer.current,
					Message:   "Invalid escape character",
				}
			}
		} else {
			builder.WriteRune(lexer.advance())
		}
	}

	if lexer.isAtEnd() {
		return &LexerError{
			Line:      lexer.line,
			Character: lexer.current,
			Message:   "Unterminated string",
		}
	}

	// Consume closing quote
	lexer.advance()

	value := builder.String()
	lexer.addToken(STRING, lexer.line, lexer.start, value)

	return nil
}

func (lexer *Lexer) unicode() (rune, *LexerError) {
	// Consume u
	lexer.advance()

	var unicodeValue int
	for i := 0; i < 4; i++ {
		if lexer.isAtEnd() {
			return 0, &LexerError{
				Line:      lexer.line,
				Character: lexer.current,
				Message:   "Unicode characters must be 4 hex digits long",
			}
		}

		nextChar := lexer.advance()
		if (nextChar < '0' || nextChar > '9') &&
			(nextChar < 'a' || nextChar > 'f') &&
			(nextChar < 'A' || nextChar > 'F') {
			return 0, &LexerError{
				Line:      lexer.line,
				Character: lexer.current,
				Message:   "Invalid unicode character",
			}
		}

		// i = 0 => math.Pow(16, 3)
		// i = 1 => math.Pow(16, 2)
		// i = 2 => math.Pow(16, 1)
		// i = 3 => math.Pow(16, 0)
		currentCharacterValue := lexer.toIntegerValue(nextChar)
		base16Power := int(math.Pow(float64(16), float64(3-i)))
		unicodeValue += currentCharacterValue * base16Power
	}

	return rune(unicodeValue), nil
}

func (lexer *Lexer) toIntegerValue(char rune) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	} else if char >= 'a' && char <= 'f' {
		return int(char - 'a' + 10)
	} else if char >= 'A' && char <= 'F' {
		return int(char - 'A' + 10)
	} else {
		// Impossible, as previously checked in unicode()
		return 0
	}
}

func (lexer *Lexer) addToken(
	tokenType TokenType,
	line int,
	character int,
	value string,
) {
	lexer.tokens = append(lexer.tokens, Token{
		Type:      tokenType,
		Line:      line,
		Character: character,
		Value:     value,
	})
}

func (lexer *Lexer) advance() rune {
	lexer.current++
	return lexer.input[lexer.current-1]
}

func (lexer *Lexer) peek() rune {
	if lexer.isAtEnd() {
		return 0
	}

	return lexer.input[lexer.current]
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.input)
}
