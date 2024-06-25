package main

import (
	"fmt"
)

type LexerError struct {
	Line      int
	Character int
	Message   string
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("Error scanning at line %d:%d - %s", e.Line, e.Character, e.Message)
}

type ParseError struct {
	Line      int
	Character int
	Message   string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Error parsing at line %d:%d - %s", e.Line, e.Character, e.Message)
}
