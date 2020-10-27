package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	if len(os.Args) != 2 {
		return errors.New("引数の個数が正しくありません")
	}

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return errors.New("引数の個数が正しくありません")
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")
	fmt.Printf("  mov rax, %d\n", i)
	fmt.Println("  ret")

	return nil
}
