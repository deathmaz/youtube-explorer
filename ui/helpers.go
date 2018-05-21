package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func regularText(v *gocui.View, text string) {
	fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", text)
}

func blueTextLn(v *gocui.View, text string, label string) {
	fmt.Fprintf(v, "\x1b[38;5;6m"+label+"%s\x1b[0m\n", text)
}

func blueText(v *gocui.View, text string, label string) {
	fmt.Fprintf(v, "\x1b[38;5;6m"+label+"%s\x1b[0m", text)
}

func highlightTextLn(v *gocui.View, text string, label string) {
	fmt.Fprintf(v, "\x1b[38;5;11m"+label+"%s\x1b[0m\n", text)
}

func headerText(v *gocui.View, text string) {
	fmt.Fprintf(v, "\x1b[33;1m%s\x1b[0m\n", text)
}
