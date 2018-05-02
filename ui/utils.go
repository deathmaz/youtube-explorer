package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// ShowLoading show loading message
func ShowLoading(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(loadingView, maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "Loading...")
		if _, err := g.SetCurrentView(loadingView); err != nil {
			return err
		}
	}

	return nil
}

// RemoveLoadin remove loading message
func RemoveLoadin(g *gocui.Gui, prevView string) error {
	if err := g.DeleteView(loadingView); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(prevView); err != nil {
		return err
	}

	return nil
}
