package main

import (
	"errors"
	"fmt"
)

var (
	ErrIncorrectNumberArgument = errors.New("the number of arguments is not correct")
	ErrNoInt                   = errors.New("this is not integer")
)

type ErrUnexpectedChar struct {
	Got  rune
	Want rune
}

func (e ErrUnexpectedChar) Error() string {
	if e.Got == '0' {
		return fmt.Sprintf("unexpected char")
	}

	if e.Want == '0' {
		return fmt.Sprintf("unexpected char %s", e.Got)
	}

	return fmt.Sprintf("unexpected char %s, but want %c", e.Got, e.Want)
}
