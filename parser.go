package main

import (
	"strconv"
	"strings"
)

// json   -> object | array | value
// object -> '{' [pair (',' pair)*] '}'
// pair   -> STRING ':' json
// array  -> '[' [json (',' json)*] ']'
// value  -> STRING | NUMBER | BOOL | NULL

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() (Node, error) {
	json, err := parser.json()
	if err != nil {
		return nil, err
	}

	return json, nil
}

// json   -> object | array | value
func (parser *Parser) json() (Node, *ParseError) {
	if parser.match(LBRACE) {
		return parser.object()
	} else if parser.match(LBRACKET) {
		return parser.array()
	} else {
		return parser.value()
	}
}

// object -> '{' [pair (',' pair)*] '}'
func (parser *Parser) object() (*Object, *ParseError) {
	if parser.match(RBRACE) {
		return &Object{
			Pairs: map[string]Node{},
		}, nil
	}

	pairs := map[string]Node{}
	// Otherwise, as matching, parse pairs
	for !parser.check(RBRACE) && !parser.isAtEnd() {
		pair, err := parser.pair()
		if err != nil {
			return nil, err
		}

		pairs[pair.Key] = pair.Value

		for parser.match(COMMA) {
			otherPair, err := parser.pair()
			if err != nil {
				return nil, err
			}

			pairs[otherPair.Key] = otherPair.Value
		}
	}

	// Encountered a right brace OR reached the end of the tokens
	if !parser.match(RBRACE) {
		return nil, parser.getError("Expect '}' after object")
	}

	return &Object{
		Pairs: pairs,
	}, nil
}

// pair   -> STRING ':' json
func (parser *Parser) pair() (*Pair, *ParseError) {
	if !parser.match(STRING) {
		return nil, parser.getError("Expected string as key")
	}

	key := parser.previous().Value

	if !parser.match(COLON) {
		return nil, parser.getError("Expected ':' after key")
	}

	value, err := parser.json()
	if err != nil {
		return nil, err
	}

	return &Pair{
		Key:   key,
		Value: value,
	}, nil
}

// array  -> '[' [json (',' json)*] ']'
func (parser *Parser) array() (*Array, *ParseError) {
	if parser.match(RBRACE) {
		return &Array{
			Elements: []Node{},
		}, nil
	}

	var elements []Node
	// Otherwise, as matching, parse pairs
	for !parser.check(RBRACKET) && !parser.isAtEnd() {
		element, err := parser.json()
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)

		for parser.match(COMMA) {
			otherElement, err := parser.json()
			if err != nil {
				return nil, err
			}

			elements = append(elements, otherElement)
		}
	}

	// Encountered a right bracket OR reached the end of the tokens
	if !parser.match(RBRACKET) {
		return nil, parser.getError("Expect ']' at the end of the array")
	}

	return &Array{
		Elements: elements,
	}, nil
}

// value  -> STRING | NUMBER | BOOL | NULL
func (parser *Parser) value() (*Value, *ParseError) {
	switch true {
	case parser.match(STRING):
		return &Value{
			Value: parser.previous().Value,
		}, nil
	case parser.match(NUMBER):
		token := parser.previous()
		//	If token value includes a dot, parse as float
		if strings.Contains(token.Value, ".") {
			parsedValue, _ := strconv.ParseFloat(token.Value, 64)
			return &Value{
				Value: parsedValue,
			}, nil
		} else {
			parsedValue, _ := strconv.Atoi(token.Value)
			return &Value{
				Value: parsedValue,
			}, nil
		}
	case parser.match(BOOL):
		return &Value{
			Value: parser.previous().Value == "true",
		}, nil
	case parser.match(NULL):
		return &Value{
			Value: nil,
		}, nil
	default:
		return nil, parser.getError("Expected value")
	}
}

func (parser *Parser) getError(message string) *ParseError {
	return &ParseError{
		Line:      parser.peek().Line,
		Character: parser.peek().Character,
		Message:   message,
	}
}

func (parser *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if parser.check(t) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) check(t TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == t
}

func (parser *Parser) peek() Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) previous() Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) isAtEnd() bool {
	return parser.current >= len(parser.tokens)
}

func (parser *Parser) advance() Token {
	if !parser.isAtEnd() {
		parser.current++
	}
	return parser.previous()
}
