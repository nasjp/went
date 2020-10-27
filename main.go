package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	numberOfArgs = 2
)

var (
	ErrIncorrectNumberArgument = errors.New("the number of arguments is not correct")
)

type ErrUnexpectedChar struct {
	Message string
}

func (e ErrUnexpectedChar) Error() string {
	return fmt.Sprintf("unexpected char, %s", e.Message)
}

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

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")

	i, n := strToInt(p)

	if n == 0 {
		return ErrUnexpectedChar{Message: p}
	}

	fmt.Printf("  mov rax, %d\n", i)

	for i := n; i < len(p); i++ {
		if p[i] == '+' {
			i++

			operand, digits := strToInt(p[i:])

			if digits == 0 {
				return ErrUnexpectedChar{Message: p}
			}

			fmt.Printf("  add rax, %d\n", operand)

			i += digits - 1

			continue
		}

		if p[i] == '-' {
			i++

			operand, digits := strToInt(p[i:])

			if digits == 0 {
				return ErrUnexpectedChar{Message: p}
			}

			fmt.Printf("  sub rax, %d\n", operand)

			i += digits - 1

			continue
		}

		return ErrUnexpectedChar{Message: p}
	}

	fmt.Println("  ret")

	return nil
}

func strToInt(s string) (int, int) {
	var (
		n         int
		digits    int = 1
		increment int
	)

	for i, c := range s {
		if c < '0' || '9' < c {
			break
		}

		n = n*digits + int(c-'0')
		digits *= 10
		increment = i + 1
	}

	return n, increment
}
