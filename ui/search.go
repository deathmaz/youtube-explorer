package ui

import (
	"net/url"
	"strings"

	"github.com/deathmaz/my-youtube/api"
	"github.com/jroimartin/gocui"
	"google.golang.org/api/youtube/v3"
)

var (
	videosList = []*youtube.SearchResult{}
)

// func customEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
// switch {
// case key == gocui.KeyEnter:
// performSearch()
// case ch != 0 && mod == 0:
// v.EditWrite(ch)
// case ch == 'j':
// v.EditWrite('j')
// case key == gocui.KeySpace:
// v.EditWrite(' ')
// case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
// v.EditDelete(true)
// }
// }

func showSearchInput(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, searchView, true); err != nil {
		return err
	}
	//  TODO: clear input //
	// view, err := g.View(searchView)
	// if err != nil {
	// return err
	// }
	// view.Clear()
	// view.MoveCursor(0, 0, false)

	deleteGlobKeybindings(g)

	return nil
}

func performSearch(g *gocui.Gui, v *gocui.View) error {
	// remove search view from history
	history = history[:len(history)-1]

	if _, err := setCurrentViewOnTop(g, searchResultsView, true); err != nil {
		return err
	}

	setGlobalKeybindings(g)
	text := v.ViewBuffer()
	if len(strings.TrimSpace(text)) == 0 {
		return nil
	}

	u, err := url.Parse(text)
	if err == nil && strings.Contains(u.Hostname(), "youtube.com") {
		vidID := u.Query().Get("v")
		// FIXME: temporary workaround of gocui paste bug
		vidID = SpaceMap(vidID)
		displayVideoPage(g, v, strings.TrimSpace(vidID))
	} else {
		viewData[searchResultsView]["query"] = text

		response, _ := api.Search(text, "video")
		videosList = response.Items
		// Group video, channel, and playlist results in separate lists.
		videos := make(map[string]string)
		/* channels := make(map[string]string)
		playlists := make(map[string]string) */

		viewData[searchResultsView]["pageToken"] = response.NextPageToken

		// Iterate through each item and add it to the correct list.
		for _, item := range response.Items {
			switch item.Id.Kind {
			case "youtube#video":
				videos[item.Id.VideoId] = item.Snippet.Title
				/* case "youtube#channel":
					channels[item.Id.ChannelId] = item.Snippet.Title
				case "youtube#playlist":
					playlists[item.Id.PlaylistId] = item.Snippet.Title */
			}
		}

		view, err := g.View(searchResultsView)
		if err != nil {
			return err
		}
		view.Clear()

		printIDs(view, "Videos", videos)
		// printIDs(v, "Channels", channels)
		// printIDs(v, "Playlists", playlists)
	}

	return nil
}

func printIDs(v *gocui.View, sectionName string, matches map[string]string) {
	// fmt.Fprintf(v, "%v:\n", sectionName)
	for id, title := range matches {
		regularText(v, "["+id+"] "+title+"")
	}
}

func goToSearchVideo(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if len(l) == 0 {
		return nil
	}

	vidID := strings.Trim(strings.Split(l, " ")[0], "[]")

	if len(vidID) == 0 {
		return nil
	}

	e := displayVideoPage(g, v, vidID)
	if e != nil {
		return e
	}

	return nil
}
