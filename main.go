package main

import (
	"fmt"
	"os"
	"strings"
)

func toObject(json string) (Node, error) {
	stripped := strings.Trim(json, "\t\n ")

	lexer := NewLexer(stripped)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	node, parseErr := parser.Parse()
	if parseErr != nil {
		return nil, err
	}

	return node, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file|input <file location | JSON input> <comma-separated list of keys>")
		os.Exit(1)
	}

	fileOrInput := os.Args[1]

	var jsonBytes []byte
	if _, err := os.Stat(fileOrInput); err == nil {
		// If file exists read input from file
		jsonBytes, err = os.ReadFile(fileOrInput)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	} else {
		// If file does not exist assume input is in JSON
		jsonBytes = []byte(fileOrInput)
	}

	json := string(jsonBytes)

	node, err := toObject(json)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Reconstructed JSON (ordering is not important):")
	fmt.Println(node)

	if len(os.Args) <= 2 {
		os.Exit(0)
	}

	keys := strings.Split(os.Args[2], ",")

	obj := ToObject(node)

	for _, key := range keys {
		if objMap, ok := obj.(map[string]interface{}); ok {
			if val, ok := objMap[key]; ok {
				obj = val
			}
		}
	}

	fmt.Println()

	fmt.Println(obj)
}
