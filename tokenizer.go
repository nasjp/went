package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type TokenKind int

const (
	TKReserved TokenKind = iota // 記号
	TKReturn                    // return
	TKIf                        // if
	TKElse                      // if
	TKFor                       // for
	TKIdent                     // 識別子
	TKNum                       // 整数
	TKEOF                       // 終点
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  int
	Str  []rune
	Loc  int
}

func NewToken(kind TokenKind, cur *Token, loc int, str ...rune) *Token {
	tok := &Token{
		Kind: kind,
		Str:  str,
		Loc:  loc,
	}

	cur.Next = tok

	return tok
}

// 次のトークンが期待値の場合には真を返す
// それ以外の場合には偽を返す
// opは省略可能.
func (tk *Token) Consume(kind TokenKind, op ...rune) bool {
	if tk.Kind != kind {
		return false
	}

	if len(op) == 0 {
		return true
	}

	if len(op) != tk.Len() {
		return false
	}

	for i := range op {
		if op[i] != tk.Str[i] {
			return false
		}
	}

	return true
}

func (tk *Token) ConsumeIdent() bool {
	return tk.Kind == TKIdent
}

// 次のトークンが期待値以外の場合にはエラーを報告する.
func (tk *Token) Expect(kind TokenKind, op ...rune) error {
	if !tk.Consume(kind, op...) {
		return userInput.Err(tk.Loc, fmt.Sprintf("'%s'ではありません", string(op)))
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

func (tk *Token) Len() int {
	return len(tk.Str)
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

		if tar := p[i:]; startsWith(tar, "==") || startsWith(tar, "!=") || startsWith(tar, "<=") || startsWith(tar, ">=") {
			cur = NewToken(TKReserved, cur, i, []rune(tar[:2])...)

			i++

			continue
		}

		if tar := p[i:]; startsWith(tar, "if") {
			cur = NewToken(TKIf, cur, i, []rune(tar[:2])...)

			i++

			continue
		}

		if tar := p[i:]; startsWith(tar, "else") {
			cur = NewToken(TKElse, cur, i, []rune(tar[:4])...)

			i += 3

			continue
		}

		if tar := p[i:]; startsWith(tar, "for") {
			cur = NewToken(TKIf, cur, i, []rune(tar[:3])...)

			i += 2

			continue
		}

		if tar := p[i:]; startsWith(tar, "return") && !isAlphaOrInt(rune(p[i+6])) {
			cur = NewToken(TKReturn, cur, i, []rune(tar[:6])...)

			i += 5

			continue
		}

		switch p[i] {
		case
			'+',
			'-',
			'*',
			'/',
			')',
			'(',
			'>',
			'<',
			'=',
			';':
			cur = NewToken(TKReserved, cur, i, []rune{rune(p[i])}...)

			continue
		}

		if isAlpha(rune(p[i])) {
			str := strToAlpha(p[i:])
			cur = NewToken(TKIdent, cur, i, []rune(str)...)

			i += len(str) - 1

			continue
		}

		if unicode.IsDigit(rune(p[i])) {
			n, err := strToInt(p[i:])
			if err != nil {
				return nil, userInput.Err(i, "数ではありません")
			}

			d := calcNumOfIntDigit(n)

			cur = NewToken(TKNum, cur, i, []rune(p[i:i+d-1])...)
			cur.Val = n

			i += d - 1

			continue
		}

		return nil, userInput.Err(i, "トークナイズできません")
	}

	NewToken(TKEOF, cur, len(p))

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
		if !isInt(c) {
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

// 文字列が対象で始まるか調べる.
func startsWith(s string, tar string) bool {
	if len(tar) > len(s) {
		return false
	}

	for i := range tar {
		if s[i] != tar[i] {
			return false
		}
	}

	return true
}

// アルファベットか調べる.
func isAlpha(s rune) bool {
	return ('a' <= s && s <= 'z') || ('A' <= s && s <= 'Z')
}

func isInt(s rune) bool {
	return ('0' <= s && s <= '9')
}

func isAlphaOrInt(s rune) bool {
	return isAlpha(s) || isInt(s) || (s == '_')
}

// 文字列をアルファベットまで読み進めるだけ読み進める.
func strToAlpha(s string) string {
	alpha := make([]rune, 0)

	for _, c := range s {
		if c < 'a' || 'z' < c {
			break
		}

		alpha = append(alpha, c)
	}

	return string(alpha)
}
