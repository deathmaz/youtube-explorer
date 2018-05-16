package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

var keyBindingsList []keyBindings

type keyBindings struct {
	view   string
	ch     rune
	key    gocui.Key
	mod    gocui.Modifier
	action func(*gocui.Gui, *gocui.View) error
}

func (k *keyBindings) getKey() interface{} {
	if k.ch == 0 {
		return k.key
	}
	return k.ch
}

func (k *keyBindings) isGlobal() bool {
	if k.view == "" {
		return true
	}

	return false
}

func keybindings(g *gocui.Gui) error {
	for _, binding := range keyBindingsList {
		if err := g.SetKeybinding(binding.view, binding.getKey(), binding.mod, binding.action); err != nil {
			log.Fatalf("Error setting keybindings: %v", err.Error())
		}
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func deleteGlobKeybindings(g *gocui.Gui) {
	for _, binding := range keyBindingsList {
		if !binding.isGlobal() {
			continue
		}
		if err := g.DeleteKeybinding(binding.view, binding.getKey(), binding.mod); err != nil {
			log.Fatalf("Error deleting keybindings: %v", err.Error())
		}
	}
}
func setGlobalKeybindings(g *gocui.Gui) {
	for _, binding := range keyBindingsList {
		if !binding.isGlobal() {
			continue
		}
		if err := g.SetKeybinding(binding.view, binding.getKey(), binding.mod, binding.action); err != nil {
			log.Fatalf("Error deleting keybindings: %v", err.Error())
		}
	}
}

func init() {
	keyBindingsList = []keyBindings{
		{
			view: "", ch: 'j', mod: gocui.ModNone, action: cursorDown,
		},
		{
			view: "", ch: 'k', mod: gocui.ModNone, action: cursorUp,
		},
		{
			view: "", ch: 'h', mod: gocui.ModNone, action: goBack,
		},
		{
			view: "", ch: 's', mod: gocui.ModNone, action: search,
		},
		{
			view: "", key: gocui.KeyArrowDown, mod: gocui.ModNone, action: cursorDown,
		},
		{
			view: "", key: gocui.KeyCtrlD, mod: gocui.ModNone, action: halfPageDown,
		},
		{
			view: "", key: gocui.KeyCtrlU, mod: gocui.ModNone, action: halfPageUp,
		},
		{
			view: "", key: gocui.KeyArrowUp, mod: gocui.ModNone, action: cursorUp,
		},
		{
			view: videoView, ch: 'p', mod: gocui.ModNone, action: playVideo,
		},
		{
			view: videoView, ch: 'g', mod: gocui.ModNone, action: selectQuality,
		},
		{
			view: videoView, ch: 'r', mod: gocui.ModNone, action: rateVideo,
		},
		{
			view: videoView, ch: 'd', mod: gocui.ModNone, action: downloadVideo,
		},
		{
			view: videosView, key: gocui.KeyEnter, mod: gocui.ModNone, action: goToVideo,
		},
		{
			view: videosView, ch: 'l', mod: gocui.ModNone, action: goToVideo,
		},
		{
			view: channelsView, key: gocui.KeyEnter, mod: gocui.ModNone, action: goToPlaylist,
		},
		{
			view: channelsView, ch: 'l', mod: gocui.ModNone, action: goToPlaylist,
		},
		{
			view: channelsView, ch: 'v', mod: gocui.ModNone, action: goToChannelVideos,
		},
		{
			view: channelPlaylistsView, key: gocui.KeyEnter, mod: gocui.ModNone, action: goToVideos,
		},
		{
			view: channelPlaylistsView, ch: 'l', mod: gocui.ModNone, action: goToVideos,
		},
		{
			view: qualityView, key: gocui.KeyEnter, mod: gocui.ModNone, action: pickQuality,
		},
		{
			view: rateVideoView, key: gocui.KeyEnter, mod: gocui.ModNone, action: rate,
		},
		{
			view: searchView, key: gocui.KeyEsc, mod: gocui.ModNone, action: goBack,
		},
	}
}
