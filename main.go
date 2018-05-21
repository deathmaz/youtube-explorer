package main

import (
	"github.com/deathmaz/my-youtube/config"
	"github.com/deathmaz/my-youtube/ui"
)

func main() {
	config.Parse()
	ui.Run()
	// api.Search("rust programming tutorials")
	// video, _ := api.GetVideos("7VcArS4Wpqk  ")
	// fmt.Println(video.Items[0].Snippet.Title)
}
