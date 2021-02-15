package plugin

import (
	"fmt"
	"strings"
)

type writer struct {
	buf strings.Builder
}

func (w *writer) p(args ...interface{}) {
	fmt.Fprint(&w.buf, args...)
	w.buf.WriteString("\n")
}

func (w *writer) a(args ...interface{}) {
	fmt.Fprint(&w.buf, args...)
}
