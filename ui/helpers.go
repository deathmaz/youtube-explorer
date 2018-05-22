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

// ShowLoading show loading message
func ShowLoading(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(loadingView, maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "Loading...")
		if _, err := setCurrentViewOnTop(g, loadingView, true); err != nil {
			return err
		}
	}

	return nil
}

// RemoveLoading remove loading message
func RemoveLoading(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(loadingView); err != nil {
		return err
	}

	goBack(g, v)
	/* if _, err := g.SetCurrentView(prevView); err != nil {
		return err
	} */

	return nil
}

func goBack(g *gocui.Gui, v *gocui.View) error {
	if len(history) > 1 {
		curView := history[len(history)-2]
		if v.Name() == searchView {
			setGlobalKeybindings(g)
		}

		history = history[:len(history)-1]

		if _, err := setCurrentViewOnTop(g, curView, false); err != nil {
			return err
		}
	}

	return nil
}
