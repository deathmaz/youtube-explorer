package ui

import (
	"fmt"
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

func showInput(g *gocui.Gui, v *gocui.View) error {
	deleteGlobKeybindings(g)

	maxX, maxY := g.Size()
	if v, err := g.SetView(searchView, maxX/2-20, maxY/2, maxX/2+20, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// var DefaultEditor gocui.Editor = gocui.EditorFunc(customEditor)

		v.Editable = true
		v.Wrap = true
		// v.Editor = DefaultEditor

		if _, err := g.SetCurrentView(searchView); err != nil {
			return err
		}
	}

	return nil
}

func performSearch(g *gocui.Gui, v *gocui.View) error {
	setGlobalKeybindings(g)
	text := v.ViewBuffer()
	if len(strings.TrimSpace(text)) > 0 {

		maxX, maxY := g.Size()
		if v, err := g.SetView(searchResultsView, 0, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Highlight = true
			v.SelBgColor = gocui.ColorYellow
			v.SelFgColor = gocui.ColorBlack
			v.Title = "Search Results"
			v.Wrap = true

			response, _ := api.Search(text, "video")
			videosList = response.Items
			// Group video, channel, and playlist results in separate lists.
			videos := make(map[string]string)
			channels := make(map[string]string)
			playlists := make(map[string]string)

			// Iterate through each item and add it to the correct list.
			for _, item := range response.Items {
				switch item.Id.Kind {
				case "youtube#video":
					videos[item.Id.VideoId] = item.Snippet.Title
				case "youtube#channel":
					channels[item.Id.ChannelId] = item.Snippet.Title
				case "youtube#playlist":
					playlists[item.Id.PlaylistId] = item.Snippet.Title
				}
			}

			printIDs(v, "Videos", videos)
			// printIDs(v, "Channels", channels)
			// printIDs(v, "Playlists", playlists)
		}

		if _, err := g.SetCurrentView(searchResultsView); err != nil {
			return err
		}
	}
	return nil
}

func printIDs(v *gocui.View, sectionName string, matches map[string]string) {
	fmt.Fprintf(v, "%v:\n", sectionName)
	for id, title := range matches {
		fmt.Fprintf(v, "[%v] %v\n", id, title)
	}
	fmt.Fprintf(v, "\n\n")
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

	maxX, maxY := g.Size()
	if v, err := g.SetView(videoView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		rating, _ := api.GetYourRating(vidID)
		commentThreads, _ := api.GetCommentThreads(vidID)
		video, _ := api.GetVideos(vidID)
		SelectedRating = rating

		if len(video.Items) > 0 {
			SelectedVideo = video.Items[0]
			fmt.Fprintf(v, "\x1b[38;5;6m%s\x1b[0m\n", video.Items[0].Id)
			fmt.Fprintf(v, "\x1b[38;5;11m%s\x1b[0m\n", video.Items[0].Snippet.PublishedAt)
			fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", video.Items[0].Snippet.Description)
			fmt.Fprintln(v, "~")
			fmt.Fprintf(v, "\x1b[38;5;6mDuration: %v\x1b[0m\n", video.Items[0].ContentDetails.Duration)
			fmt.Fprintf(v, "\x1b[38;5;6mTotal views: %v\x1b[0m\n", video.Items[0].Statistics.ViewCount)
			fmt.Fprintf(v, "\x1b[38;5;6mLikes: %v\x1b[0m\n", video.Items[0].Statistics.LikeCount)
			fmt.Fprintf(v, "\x1b[38;5;6mDislikes: %v\x1b[0m\n~", video.Items[0].Statistics.DislikeCount)
		}

		fmt.Fprintln(v, "")
		fmt.Fprintf(v, "\x1b[38;5;208mYour rating: %s\x1b[0m\n~\n~\n", rating)
		fmt.Fprint(v, "\x1b[33;1mComments:\x1b[0m\n~\n")

		for _, thread := range commentThreads.Items {
			comment := thread.Snippet.TopLevelComment
			fmt.Fprintf(v, "\x1b[38;5;6m%s\x1b[0m ", comment.Snippet.AuthorDisplayName)
			fmt.Fprintf(v, "\x1b[38;5;6m%v %s\x1b[0m \n", comment.Snippet.LikeCount, "Likes")
			fmt.Fprintf(v, "\x1b[38;5;11m%s\x1b[0m\n", comment.Snippet.PublishedAt)
			fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", comment.Snippet.TextDisplay)

			if thread.Replies != nil {
				fmt.Fprint(v, "\x1b[33;1mReplies:\x1b[0m")
				comments := thread.Replies.Comments
				for i := len(comments) - 1; i >= 0; i-- {
					fmt.Fprintf(v, "\n    \x1b[38;5;6m%s\x1b[0m ", comments[i].Snippet.AuthorDisplayName)
					fmt.Fprintf(v, "    \x1b[38;5;6m%v %s\x1b[0m \n", comments[i].Snippet.LikeCount, "Likes")
					fmt.Fprintf(v, "    \x1b[38;5;11m%s\x1b[0m\n", comments[i].Snippet.PublishedAt)
					fmt.Fprintf(v, "    \x1b[38;5;3m%s\x1b[0m\n", comments[i].Snippet.TextDisplay)
				}
			}
			fmt.Fprintln(v, "")
		}

		if _, err := g.SetCurrentView(videoView); err != nil {
			return err
		}
	}

	return nil
}
