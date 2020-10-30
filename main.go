package main

import (
	"fmt"
	"os"
)

const (
	numberOfArgs = 2
)

// 現在着目しているトークン.
var currentToken *Token

func proceedToken() {
	currentToken = currentToken.Next
}

// ユーザーの入力文字列を保持する.
var userInput UserInput

// ニーモニックのラベル名を管理する.
var label int

// local変数保存用
// 関数ごとに初期化される.
var localValue *Node

// 出力先
var output *Writer

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

	output = NewWriter(os.Stdout)

	p := os.Args[1]

	userInput = UserInput(p)

	token, err := tokenize(p)
	if err != nil {
		return err
	}

	currentToken = token

	node, err := parse()
	if err != nil {
		return err
	}

	if err := generate(node); err != nil {
		return err
	}

	return nil
}
