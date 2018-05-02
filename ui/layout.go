package ui

import (
	"fmt"

	"github.com/deathmaz/my-youtube/api"
	"github.com/jroimartin/gocui"
)

// Layout layout setup
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(channelsView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = channelsView
		v.Wrap = true

		if _, err = setCurrentViewOnTop(g, channelsView); err != nil {
			return err
		}

		if len(subscriptions) == 0 {
			go func() {
				ShowLoading(g)
				response, _ := api.GetMySubscriptions()
				subscriptions = response.Items

				g.Update(func(g *gocui.Gui) error {
					v, err := g.View(channelsView)
					if err != nil {
						return err
					}
					v.Clear()

					for _, channel := range subscriptions {
						fmt.Fprintln(v, channel.Snippet.Title)
					}
					RemoveLoadin(g, v.Title)
					return nil
				})
			}()
		} else {
			for _, channel := range subscriptions {
				fmt.Fprintln(v, channel.Snippet.Title)
			}
		}

	}

	/* if v, err := g.SetView("v3", 0, maxY/2-1, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3"
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprint(v, "Press TAB to change current view")
	}
	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v4 (editable)"
		v.Editable = true
	} */
	return nil
}
