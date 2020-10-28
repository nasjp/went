package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrIncorrectNumberArgument = errors.New("the number of arguments is not correct")
	ErrNoInt                   = errors.New("this is not integer")
)

type UserInput string

func (ui UserInput) Err(loc int, message string) error {
	body := fmt.Sprintf(`%s
%s^ %s`, ui, strings.Repeat(" ", loc), message)

	return InvalidInputError{s: body}
}

type InvalidInputError struct{ s string }

func (e InvalidInputError) Error() string { return e.s }
