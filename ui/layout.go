package ui

import (
	"fmt"

	"github.com/deathmaz/my-youtube/api"
	"github.com/jroimartin/gocui"
)

var (
	history = []string{}
)

func setCurrentViewOnTop(g *gocui.Gui, name string, writeHistory bool) (*gocui.View, error) {
	if writeHistory {
		history = append(history, name)
	}
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(channelPlaylistsView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Channels playlists view"
		v.Wrap = true
	}

	if v, err := g.SetView(searchResultsView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Search Results"
		v.Wrap = true
	}

	if v, err := g.SetView(searchView, maxX/2-20, maxY/2, maxX/2+20, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// var DefaultEditor gocui.Editor = gocui.EditorFunc(customEditor)

		v.Editable = true
		v.Wrap = true
		// v.Editor = DefaultEditor
	}

	if v, err := g.SetView(rateVideoView, maxX/2-15, maxY/2-3, maxX/2+15, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Rate Video"
	}

	if v, err := g.SetView(qualityView, maxX/2-15, maxY/2-3, maxX/2+15, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Select Video quality"
	}

	if v, err := g.SetView(videoView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
	}

	if v, err := g.SetView(videosView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Videos"
		v.Wrap = true
	}

	if v, err := g.SetView(channelsView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Title = channelsView
		v.Wrap = true
		if len(subscriptions) == 0 {
			ShowLoading(g)
			response, _ := api.GetMySubscriptions()
			subscriptions = response.Items

			v, err := g.View(channelsView)
			if err != nil {
				return err
			}
			v.Clear()

			for _, channel := range subscriptions {
				fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", channel.Snippet.Title)
			}
			RemoveLoading(g, v.Title)
		} else {
			for _, channel := range subscriptions {
				fmt.Fprintln(v, channel.Snippet.Title)
			}
		}

		if _, err := setCurrentViewOnTop(g, channelsView, true); err != nil {
			return err
		}

	}

	return nil
}
