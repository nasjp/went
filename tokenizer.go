package main

import (
	"fmt"
	"strconv"
	"unicode"
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
	Loc  int
}

func NewToken(kind TokenKind, cur *Token, str []rune, loc int) *Token {
	tok := &Token{
		Kind: kind,
		Str:  str,
		Loc:  loc,
	}

	cur.Next = tok

	return tok
}

// 次のトークンが期待している記号のときには真を返す
// それ以外の場合には偽を返す.
func (tk *Token) Consume(op rune) bool {
	if tk.Kind != TKReserved || tk.Str[0] != op {
		return false
	}

	return true
}

// 次のトークンが期待値以外の場合にはエラーを報告する.
func (tk *Token) Expect(op rune) error {
	if tk.Kind != TKReserved || tk.Str[0] != op {
		return userInput.Err(tk.Loc, fmt.Sprintf("'%c'ではありません", op))
	}

	return nil
}

// 次のトークンが数値の場合、トークンを1つ読み進めてその数値を返す
// それ以外の場合にはエラーを報告する.
func (tk *Token) ExpectNum() (int, error) {
	if tk.Kind != TKNum {
		return 0, userInput.Err(tk.Loc, "数ではありません")
	}

	val := tk.Val

	return val, nil
}

func (tk *Token) AtEOF() bool {
	return tk.Kind == TKEOF
}

func tokenize(p string) (*Token, error) {
	head := &Token{}
	cur := head

	for i := 0; i < len(p); i++ {
		if unicode.IsSpace(rune(p[i])) {
			continue
		}

		switch p[i] {
		case
			'+',
			'-',
			'*',
			'/',
			')',
			'(':
			cur = NewToken(TKReserved, cur, []rune{rune(p[i])}, i)

			continue
		}

		if unicode.IsDigit(rune(p[i])) {
			n, err := strToInt(p[i:])
			if err != nil {
				return nil, userInput.Err(cur.Loc, "数ではありません")
			}

			d := calcNumOfIntDigit(n)

			cur = NewToken(TKNum, cur, []rune(p[i:i+d-1]), i)
			cur.Val = n

			i += d - 1

			continue
		}

		return nil, userInput.Err(cur.Loc, "トークナイズできません")
	}

	NewToken(TKEOF, cur, nil, len(p))

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
