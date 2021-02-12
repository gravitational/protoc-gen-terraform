package plugin

import (
	"fmt"
	"strings"
)

type writer struct {
	buf strings.Builder
	//	nesting int
}

func (w *writer) p(args ...interface{}) {
	fmt.Fprint(&w.buf, args...)
	w.buf.WriteString("\n")
}

func (w *writer) a(args ...interface{}) {
	fmt.Fprint(&w.buf, args...)
}

// func (w *writer) in() {
// 	w.nesting++
// 	w.buf.WriteString("{")
// }

// func (w *writer) out() {
// 	for i := 0; i < w.nesting; i++ {
// 		w.buf.WriteString("}")
// 	}
// 	w.buf.WriteString("\n")
// 	w.nesting = 0
// }
