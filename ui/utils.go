package ui

import (
	"fmt"
	"math"

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

func goBack(g *gocui.Gui, v *gocui.View) error {
	var views []*gocui.View
	for _, view := range g.Views() {
		if view.Name() != loadingView {
			views = append(views, view)
		}
	}

	if len(views) > 1 {
		if v.Name() == searchView {
			setGlobalKeybindings(g)
		}

		if err := g.DeleteView(views[len(views)-1].Name()); err != nil {
			return err
		}

		if _, err := g.SetCurrentView(views[len(views)-2].Name()); err != nil {
			return err
		}

	}
	return nil
}

// Round round
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
