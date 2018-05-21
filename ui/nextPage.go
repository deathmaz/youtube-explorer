package ui

import (
	"fmt"

	"github.com/deathmaz/my-youtube/api"
	"github.com/jroimartin/gocui"
)

func nextPage(g *gocui.Gui, v *gocui.View) error {
	nextPageToken, ok := viewData[v.Name()]["pageToken"]

	if ok && len(nextPageToken) == 0 {
		return nil
	}

	switch v.Name() {
	case channelsView:
		res, _ := api.GetMySubscriptionsNextPage(nextPageToken)
		if len(res.Items) > 0 {
			viewData[v.Name()]["pageToken"] = res.NextPageToken

			for _, channel := range res.Items {
				subscriptions = append(subscriptions, channel)
				fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", channel.Snippet.Title)
			}
		}

	case videosView:
		playlistID := viewData[v.Name()]["playlistID"]
		res, _ := api.GetPlaylistItemsNextPage(playlistID, nextPageToken)
		if len(res.Items) > 0 {
			viewData[v.Name()]["pageToken"] = res.NextPageToken

			for _, video := range res.Items {
				videos = append(videos, video)
				fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", video.Snippet.Title)
			}
		}

	case searchResultsView:
		res, _ := api.GetSearchNextPage(viewData[v.Name()]["query"], "video", viewData[v.Name()]["pageToken"])
		viewData[searchResultsView]["pageToken"] = res.NextPageToken
		videos := make(map[string]string)

		for _, item := range res.Items {
			switch item.Id.Kind {
			case "youtube#video":
				videosList = append(videosList, item)
				videos[item.Id.VideoId] = item.Snippet.Title
			}
		}
		printIDs(v, "Videos", videos)

	case channelPlaylistsView:
		res, _ := api.GetChannelPlaylistItemsNextPage(viewData[channelPlaylistsView]["channelID"], viewData[channelPlaylistsView]["pageToken"])
		viewData[channelPlaylistsView]["pageToken"] = res.NextPageToken
		for _, playlist := range res.Items {
			playlists = append(playlists, playlist)
			fmt.Fprintf(v, "\x1b[38;5;3m%s\x1b[0m\n", playlist.Snippet.Title)
		}
	}

	return nil
}
