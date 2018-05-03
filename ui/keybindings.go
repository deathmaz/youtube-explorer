package ui

import "github.com/jroimartin/gocui"

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(videoView, 'p', gocui.ModNone, playVideo); err != nil {
		return err
	}

	if err := g.SetKeybinding(videoView, 'g', gocui.ModNone, selectQuality); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'h', gocui.ModNone, goBack); err != nil {
		return err
	}

	if err := g.SetKeybinding(qualityView, gocui.KeyEnter, gocui.ModNone, pickQuality); err != nil {
		return err
	}

	if err := g.SetKeybinding(videoView, 'r', gocui.ModNone, rateVideo); err != nil {
		return err
	}

	if err := g.SetKeybinding(videoView, 'd', gocui.ModNone, downloadVideo); err != nil {
		return err
	}

	if err := g.SetKeybinding(rateVideoView, gocui.KeyEnter, gocui.ModNone, rate); err != nil {
		return err
	}

	if err := g.SetKeybinding(videosView, gocui.KeyEnter, gocui.ModNone, goToVideo); err != nil {
		return err
	}

	if err := g.SetKeybinding(channelsView, gocui.KeyEnter, gocui.ModNone, goToPlaylist); err != nil {
		return err
	}

	if err := g.SetKeybinding(channelsView, 'v', gocui.ModNone, goToChannelVideos); err != nil {
		return err
	}

	if err := g.SetKeybinding(channelPlaylistsView, gocui.KeyEnter, gocui.ModNone, goToVideos); err != nil {
		return err
	}

	if err := g.SetKeybinding(msgView, gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	return nil
}
