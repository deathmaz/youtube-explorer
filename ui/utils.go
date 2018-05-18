package ui

import (
	"fmt"
	"log"
	"math"
	"os/exec"

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

// RemoveLoading remove loading message
func RemoveLoading(g *gocui.Gui, prevView string) error {
	if err := g.DeleteView(loadingView); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(prevView); err != nil {
		return err
	}

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

func runcmd(cmd string, shell bool) []byte {
	if shell {
		err := exec.Command("bash", "-c", cmd).Start()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
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
