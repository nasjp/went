package main

const offsetSize = 8

type LocalValue struct {
	Next   *LocalValue
	Name   []rune
	Offset int
}

func NewLocalValue(tk *Token, beforeOffset int) *LocalValue {
	lv := &LocalValue{
		Name:   tk.Str,
		Offset: beforeOffset + offsetSize,
	}

	lv.Next = localValue

	return lv
}

func (lv *LocalValue) Len() int {
	return len(lv.Name)
}

func findLocalValue(tok *Token) *LocalValue {
	for val := localValue; val != nil; val = val.Next {
		if val.Len() == tok.Len() && string(tok.Str) == string(val.Name) {
			return val
		}
	}

	return nil
}
