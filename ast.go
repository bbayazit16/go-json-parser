package main

import (
	"fmt"
	"strings"
)

// json   -> object | array | value
// object -> '{' [pair (',' pair)*] '}'
// pair   -> STRING ':' json
// array  -> '[' [json (',' json)*] ']'
// value  -> STRING | NUMBER | BOOL | NULL

type Node interface {
	String() string
}

type Object struct {
	Pairs map[string]Node
}

type Pair struct {
	Key   string
	Value Node
}

type Array struct {
	Elements []Node
}

type Value struct {
	Value interface{}
}

func (o Object) String() string {
	pairs := make([]string, 0, len(o.Pairs))
	for k, v := range o.Pairs {
		pairs = append(pairs, fmt.Sprintf("\"%s\": %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (p Pair) String() string {
	return fmt.Sprintf("\"%s\": %s", p.Key, p.Value.String())
}

func (a Array) String() string {
	elements := make([]string, len(a.Elements))
	for i, el := range a.Elements {
		elements[i] = el.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (v Value) String() string {
	switch val := v.Value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case float64, int, bool:
		return fmt.Sprintf("%v", val)
	case nil:
		return "null"
	default:
		return ""
	}
}
