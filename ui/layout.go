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
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = channelsView
		v.Wrap = true

		if _, err := g.SetCurrentView(channelsView); err != nil {
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
						fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", channel.Snippet.Title)
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
	return nil
}
