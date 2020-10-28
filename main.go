package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

const (
	numberOfArgs = 2
)

type TokenKind int

const (
	TKReserved TokenKind = iota
	TKNum
	TKEOF
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  int
	Str  []rune
}

func NewToken(kind TokenKind, cur *Token, str []rune) *Token {
	tok := &Token{
		Kind: kind,
		Str:  str,
	}

	cur.Next = tok

	return tok
}

// 現在着目しているトークン
var token *Token

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

	var err error
	if token, err = tokenize(p); err != nil {
		return err
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")

	n, err := expectNum()
	if err != nil {
		return err
	}

	fmt.Printf("  mov rax, %d\n", n)

	for !atEOF() {
		if consume('+') {
			operand, err := expectNum()
			if err != nil {
				return err
			}

			fmt.Printf("  add rax, %d\n", operand)

			continue
		}

		if consume('-') {
			operand, err := expectNum()
			if err != nil {
				return err
			}

			fmt.Printf("  sub rax, %d\n", operand)

			continue
		}

		return ErrUnexpectedChar{}
	}

	fmt.Println("  ret")

	return nil
}

func tokenize(p string) (*Token, error) {
	head := &Token{}
	cur := head

	for i := 0; i < len(p); i++ {
		if unicode.IsSpace(rune(p[i])) {
			continue
		}

		if p[i] == '+' || p[i] == '-' {
			cur = NewToken(TKReserved, cur, []rune{rune(p[i])})

			continue
		}

		if unicode.IsDigit(rune(p[i])) {
			n, err := strToInt(p[i:])
			if err != nil {
				return nil, err
			}

			d := calcNumOfIntDigit(n)

			cur = NewToken(TKNum, cur, []rune(p[i:i+d-1]))
			cur.Val = n

			i += d - 1

			continue
		}

		return nil, ErrUnexpectedChar{}
	}

	NewToken(TKEOF, cur, nil)

	return head.Next, nil
}

// 文字列を整数値まで読み進めるだけ読み進める.
func strToInt(s string) (int, error) {
	var (
		n      int
		digits int = 1
		read   bool
	)

	for _, c := range s {
		if c < '0' || '9' < c {
			break
		}

		n = n*digits + int(c-'0')
		digits *= 10
		read = true
	}

	if !read {
		return 0, ErrNoInt
	}

	return n, nil
}

// 整数値の桁数を調べる.
func calcNumOfIntDigit(n int) int {
	return len(strconv.Itoa(n))
}

// 次のトークンが期待している記号のときには、トークンを1つ読み進めて
// 真を返す。それ以外の場合には偽を返す.
func consume(op rune) bool {
	if token.Kind != TKReserved || token.Str[0] != op {
		return false
	}

	token = token.Next

	return true
}

// 次のトークンが数値の場合、トークンを1つ読み進めてその数値を返す。
// それ以外の場合にはエラーを報告する.
func expect(op rune) error {
	if token.Kind != TKReserved || token.Str[0] != op {
		return ErrUnexpectedChar{Want: op, Got: token.Str[0]}
	}

	token = token.Next

	return nil
}

func expectNum() (int, error) {
	if token.Kind != TKNum {
		return 0, ErrNoInt
	}

	val := token.Val
	token = token.Next

	return val, nil
}

func atEOF() bool {
	return token.Kind == TKEOF
}
