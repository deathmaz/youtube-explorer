package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deathmaz/my-youtube/api"
	"github.com/deathmaz/my-youtube/config"
	"github.com/jroimartin/gocui"
	"google.golang.org/api/youtube/v3"
)

var (
	subscriptions        = []*youtube.Subscription{}
	playlists            = []*youtube.Playlist{}
	videos               = []*youtube.PlaylistItem{}
	selectedVideoQuality = "720"
	videoQuality         = []string{"360", "480", "720", "1080"}
	ratings              = []string{"like", "dislike", "none"}

	// SelectedRating selected rating for video
	SelectedRating = ""
	// SelectedVideo selected video
	SelectedVideo *youtube.Video
	viewData      = make(map[string]map[string]string)
)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	bufLines := v.BufferLines()
	if v != nil && len(bufLines) > 0 {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		maxX, _ := g.Size()
		lines := 0
		for _, line := range bufLines[1 : len(bufLines)-1] {
			if len(line) > maxX-2 {
				lines += int(Round(float64(len(line)/maxX), .1, 0))
			} else {
				lines++
			}
		}

		if oy+cy < lines {
			if err := v.SetCursor(cx, cy+1); err != nil {
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return nil
}

func halfPageDown(g *gocui.Gui, v *gocui.View) error {
	bufLines := v.BufferLines()
	if v != nil && len(bufLines) > 0 {
		cx, cy := v.Cursor()
		maxX, maxY := g.Size()
		ox, oy := v.Origin()

		curY := cy + maxY/2
		lines := 0
		for _, line := range bufLines[1 : len(bufLines)-1] {
			if len(line) > maxX-2 {
				lines += int(Round(float64(len(line)/maxX), .1, 0))
			} else {
				lines++
			}
		}

		if oy+curY > lines {
			curY = lines - oy
		}

		if err := v.SetCursor(cx, curY); err != nil {
			if err := v.SetOrigin(ox, oy+maxY/2); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func halfPageUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		_, maxY := g.Size()
		cursorMaxY := cy - maxY/2
		originMaxY := oy - maxY/2
		if oy <= 0 {
			cursorMaxY = 0
		}

		if originMaxY < 0 {
			originMaxY = 0
		}
		if err := v.SetCursor(cx, cursorMaxY); err != nil {
			if err := v.SetOrigin(ox, originMaxY); err != nil {
				return err
			}
		}
	}
	return nil
}

func goToPlaylists(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := setCurrentViewOnTop(g, channelPlaylistsView, true); err != nil {
		return err
	}

	view, err := g.View(channelPlaylistsView)
	if err != nil {
		return err
	}
	view.Clear()

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	for _, subscription := range subscriptions {
		if subscription.Snippet.Title == l {
			res, _ := api.ChannelPlaylistItems(subscription.Snippet.ResourceId.ChannelId)
			playlists = res.Items
			viewData[channelPlaylistsView]["pageToken"] = res.NextPageToken
			viewData[channelPlaylistsView]["channelID"] = subscription.Snippet.ResourceId.ChannelId

			for _, playlist := range playlists {
				regularText(view, playlist.Snippet.Title)
			}
			break
		}
	}

	return nil
}

func goToVideos(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if len(l) == 0 {
		return nil
	}

	view, err := g.View(videosView)
	if err != nil {
		return err
	}
	view.Clear()

	if _, err := setCurrentViewOnTop(g, videosView, true); err != nil {
		return err
	}

	for _, playlist := range playlists {
		if playlist.Snippet.Title == l {
			res, _ := api.PlaylistItems(playlist.Id)
			videos = res.Items
			viewData[videosView]["pageToken"] = res.NextPageToken
			viewData[videosView]["playlistID"] = playlist.Id

			for _, video := range videos {
				regularText(view, video.Snippet.Title)
			}

			break
		}
	}

	return nil

}

func goToVideoChannelPlaylists(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, channelPlaylistsView, true); err != nil {
		return err
	}

	view, err := g.View(channelPlaylistsView)
	if err != nil {
		return err
	}
	view.Clear()

	res, _ := api.ChannelPlaylistItems(SelectedVideo.Snippet.ChannelId)
	playlists = res.Items
	viewData[channelPlaylistsView]["pageToken"] = res.NextPageToken
	viewData[channelPlaylistsView]["channelID"] = SelectedVideo.Snippet.ChannelId

	for _, playlist := range playlists {
		regularText(view, playlist.Snippet.Title)
	}

	return nil
}

func goToVideoChannelVideos(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, videosView, true); err != nil {
		return err
	}

	view, err := g.View(videosView)
	if err != nil {
		return err
	}

	view.Clear()

	channel, _ := api.Channel(SelectedVideo.Snippet.ChannelId)
	res, _ := api.PlaylistItems(channel.Items[0].ContentDetails.RelatedPlaylists.Uploads)
	videos = res.Items
	viewData[videosView]["pageToken"] = res.NextPageToken
	viewData[videosView]["playlistID"] = channel.Items[0].ContentDetails.RelatedPlaylists.Uploads

	for _, video := range videos {
		regularText(view, video.Snippet.Title)
	}

	return nil
}

func goToChannelVideos(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	view, err := g.View(videosView)
	if err != nil {
		return err
	}

	view.Clear()

	for _, subscription := range subscriptions {
		if subscription.Snippet.Title == l {
			channel, _ := api.Channel(subscription.Snippet.ResourceId.ChannelId)
			res, _ := api.PlaylistItems(channel.Items[0].ContentDetails.RelatedPlaylists.Uploads)
			videos = res.Items
			viewData[videosView]["pageToken"] = res.NextPageToken
			viewData[videosView]["playlistID"] = channel.Items[0].ContentDetails.RelatedPlaylists.Uploads

			for _, video := range videos {
				regularText(view, video.Snippet.Title)
			}

			if _, err := setCurrentViewOnTop(g, videosView, true); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func goToVideo(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	for _, video := range videos {
		if video.Snippet.Title == l {
			e := displayVideoPage(g, v, video.ContentDetails.VideoId)
			if e != nil {
				return e
			}
			break
		}

	}

	return nil
}

func displayVideoPage(g *gocui.Gui, v *gocui.View, vidID string) error {
	if _, err := setCurrentViewOnTop(g, videoView, true); err != nil {
		return err
	}

	view, err := g.View(videoView)
	view.Clear()
	if err != nil {
		return err
	}

	video, _ := api.Videos(vidID)

	if len(video.Items) > 0 {
		rating, _ := api.YourRating(vidID)
		SelectedRating = rating
		commentThreads, _ := api.CommentThreads(vidID)

		vid := video.Items[0]
		view.Title = vid.Snippet.Title
		SelectedVideo = vid
		blueTextLn(view, vid.Id, "")
		highlightTextLn(view, vid.Snippet.PublishedAt, "")
		regularText(view, vid.Snippet.Description)
		fmt.Fprintln(view, "")
		blueTextLn(view, vid.ContentDetails.Duration, "Duration: ")
		blueTextLn(view, strconv.FormatUint(vid.Statistics.ViewCount, 10), "Total view: ")
		blueTextLn(view, strconv.FormatUint(vid.Statistics.LikeCount, 10), "Likes: ")
		blueTextLn(view, strconv.FormatUint(vid.Statistics.DislikeCount, 10), "Dislikes: ")

		fmt.Fprintln(view, "")
		fmt.Fprintf(view, "\x1b[38;5;208mYour rating: %s\x1b[0m\n~\n~\n", rating)
		headerText(view, "Comments:")

		for _, thread := range commentThreads.Items {
			comment := thread.Snippet.TopLevelComment
			blueText(view, comment.Snippet.AuthorDisplayName+" ", "")
			blueTextLn(view, strconv.FormatInt(comment.Snippet.LikeCount, 10)+" Likes", "")
			highlightTextLn(view, comment.Snippet.PublishedAt, "")
			regularText(view, comment.Snippet.TextDisplay)

			if thread.Replies != nil {
				headerText(view, "Replies:")
				comments := thread.Replies.Comments
				for i := len(comments) - 1; i >= 0; i-- {
					blueText(view, "    "+comments[i].Snippet.AuthorDisplayName+" ", "")
					blueTextLn(view, strconv.FormatInt(comments[i].Snippet.LikeCount, 10)+" Likes", "")
					highlightTextLn(view, "    "+comments[i].Snippet.PublishedAt, "")
					fmt.Fprint(view, "    ")
					regularText(view, comments[i].Snippet.TextDisplay)
				}
			}
			fmt.Fprintln(view, "")
		}
	} else {
		regularText(view, "No result")
	}

	return nil
}

func playVideo(g *gocui.Gui, v *gocui.View) error {
	cmd := exec.Command("bash", "-c",
		"mpv --ytdl-format='bestvideo[height<="+selectedVideoQuality+"]+bestaudio/best[height<="+selectedVideoQuality+"]' https://www.youtube.com/watch?v="+SelectedVideo.Id)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func downloadVideo(g *gocui.Gui, v *gocui.View) error {
	tspCommand := ""
	usr, _ := user.Current()
	outPath := filepath.FromSlash(usr.HomeDir + "/Downloads/")
	if _, err := os.Stat(config.Conf.DownloadPath); err == nil {
		if strings.HasSuffix(config.Conf.DownloadPath, "/") {
			outPath = config.Conf.DownloadPath
		} else {
			outPath = config.Conf.DownloadPath + "/"
		}
	}

	if path, err := exec.LookPath("tsp"); err == nil {
		tspCommand += "TS_SOCKET=/tmp/y-dl " + path + " "
	}

	cmd := exec.Command("bash", "-c",
		tspCommand+"youtube-dl -o '"+outPath+"%(title)s.%(ext)s' -f 'bestvideo[height<="+selectedVideoQuality+"]+bestaudio/best[height<="+selectedVideoQuality+"]' https://www.youtube.com/watch?v="+SelectedVideo.Id)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func selectQuality(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, qualityView, true); err != nil {
		return err
	}

	view, err := g.View(qualityView)
	if err != nil {
		return err
	}
	view.Clear()

	for _, quality := range videoQuality {
		if quality == selectedVideoQuality {
			highlightTextLn(view, quality, "")
		} else {
			fmt.Fprintln(view, quality)
		}
	}

	return nil
}

func pickQuality(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	selectedVideoQuality = l
	goBack(g, v)

	return nil
}

func rateVideo(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, rateVideoView, true); err != nil {
		return err
	}

	view, err := g.View(rateVideoView)
	if err != nil {
		return err
	}
	view.Clear()

	for _, rating := range ratings {
		if SelectedRating == rating {
			highlightTextLn(view, rating, "")
		} else {
			fmt.Fprintln(view, rating)
		}
	}

	return nil
}

func rate(g *gocui.Gui, v *gocui.View) error {
	if strings.HasPrefix(v.Name(), rateVideoView) {
		var l string
		var err error

		_, cy := v.Cursor()
		if l, err = v.Line(cy); err != nil {
			l = ""
		}

		e := api.RateVideo(SelectedVideo.Id, l)
		if e == nil {
			goBack(g, v)
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(msgView); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(channelsView); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Run func
func Run() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.InputEsc = true
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init() {
	viewData[videosView] = map[string]string{
		"pageToken":  "",
		"playlistID": "",
	}

	viewData[channelsView] = map[string]string{
		"pageToken": "",
	}

	viewData[channelPlaylistsView] = map[string]string{
		"pageToken": "",
		"channelID": "",
	}

	viewData[searchResultsView] = map[string]string{
		"pageToken": "",
		"query":     "",
	}
}
