package main

import (
	"fmt"
	"io"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) F(format string, a ...interface{}) {
	fmt.Fprintf(w.w, format, a...)
}

func (w *Writer) L(a ...interface{}) {
	fmt.Fprintln(w.w, a...)
}
