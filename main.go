package main

import (
	"fmt"
	"os"
)

const (
	numberOfArgs = 2
)

// 現在着目しているトークン.
var token *Token

func proceedToken() {
	token = token.Next
}

// ユーザーの入力文字列を保持する.
var userInput UserInput

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	if len(os.Args) != numberOfArgs {
		return ErrIncorrectNumberArgument
	}

	p := os.Args[1]

	userInput = UserInput(p)

	var err error
	if token, err = tokenize(p); err != nil {
		return err
	}

	node, err := expr()
	if err != nil {
		return err
	}

	prequel()

	generate(node)

	sequel()

	return nil
}
