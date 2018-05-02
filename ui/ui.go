package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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
	selectedVideo        *youtube.PlaylistItem
	selectedVideoQuality = "720"
	videoQuality         = []string{"360", "480", "720", "1080"}
	ratings              = []string{"like", "dislike", "none"}
)

func runcmd(cmd string, shell bool) []byte {
	if shell {
		err := exec.Command("bash", "-c", cmd).Start()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func goToChannels(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(channelPlaylistsView); err != nil {
		return err
	}

	playlists = []*youtube.Playlist{}

	if _, err := g.SetCurrentView(channelsView); err != nil {
		return err
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

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func delVideoView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(videoView); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(videosView); err != nil {
		return err
	}

	return nil
}

func delVideosView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(videosView); err != nil {
		return err
	}

	if len(playlists) > 0 {
		if _, err := g.SetCurrentView(channelPlaylistsView); err != nil {
			return err
		}
	} else {
		if _, err := g.SetCurrentView(channelsView); err != nil {
			return err
		}
	}

	return nil
}

func goToPlaylist(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView(channelPlaylistsView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Channels playlists view"
		v.Wrap = true

		for _, subscription := range subscriptions {
			if subscription.Snippet.Title == l {
				go func() {
					ShowLoading(g)
					res, _ := api.GetChannelPlaylistItems(subscription.Snippet.ResourceId.ChannelId)
					playlists = res.Items

					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(channelPlaylistsView)
						if err != nil {
							return err
						}
						v.Clear()

						for _, playlist := range playlists {
							fmt.Fprintln(v, playlist.Snippet.Title)
						}
						RemoveLoadin(g, channelPlaylistsView)
						return nil
					})
				}()
				break
			}
		}

		if _, err := g.SetCurrentView(channelPlaylistsView); err != nil {
			return err
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

	maxX, maxY := g.Size()
	if v, err := g.SetView(videosView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Videos"
		v.Wrap = true

		for _, playlist := range playlists {
			if playlist.Snippet.Title == l {
				go func() {
					ShowLoading(g)
					res, _ := api.GetPlaylistItems(playlist.Id)
					videos = res.Items

					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(videosView)
						if err != nil {
							return err
						}
						v.Clear()

						for _, video := range videos {
							fmt.Fprintln(v, video.Snippet.Title)
						}
						RemoveLoadin(g, videosView)
						return nil
					})
				}()
				break
			}
		}

		if _, err := g.SetCurrentView(channelPlaylistsView); err != nil {
			return err
		}
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

	maxX, maxY := g.Size()
	if v, err := g.SetView(videosView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Videos"
		v.Wrap = true

		for _, subscription := range subscriptions {
			if subscription.Snippet.Title == l {
				go func() {
					ShowLoading(g)
					channel, _ := api.GetChannel(subscription.Snippet.ResourceId.ChannelId)
					res, _ := api.GetPlaylistItems(channel.Items[0].ContentDetails.RelatedPlaylists.Uploads)
					videos = res.Items

					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(videosView)
						if err != nil {
							return err
						}
						v.Clear()

						for _, video := range videos {
							fmt.Fprintln(v, video.Snippet.Title)
						}
						RemoveLoadin(g, videosView)
						return nil
					})
				}()
				break
			}
		}

		if _, err := g.SetCurrentView(videosView); err != nil {
			return err
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

	maxX, maxY := g.Size()
	if v, err := g.SetView(videoView, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Wrap = true

		fmt.Fprintln(v, l)
		for _, video := range videos {
			if video.Snippet.Title == l {
				selectedVideo = video

				v.Title = video.Snippet.Title
				fmt.Fprintln(v, video.ContentDetails.VideoId)
				fmt.Fprintln(v, video.ContentDetails.VideoPublishedAt)
				fmt.Fprintln(v, video.Snippet.Description)

				go func() {
					rating, _ := api.GetRating(video.ContentDetails.VideoId)
					comments, _ := api.GetComments(video.ContentDetails.VideoId)

					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(videoView)
						if err != nil {
							return err
						}

						fmt.Fprintln(v, "")
						fmt.Fprintln(v, "Your rating:", rating)
						fmt.Fprintln(v, "")
						fmt.Fprintln(v, "")

						for _, comment := range comments.Items {
							fmt.Fprintln(v, comment.Snippet.AuthorDisplayName, comment.Snippet.LikeCount, "Likes")
							fmt.Fprintln(v, comment.Snippet.PublishedAt)
							fmt.Fprintln(v, comment.Snippet.TextDisplay)
							fmt.Fprintln(v, "")
						}

						return nil
					})
				}()

				break
			}
		}

		if _, err := g.SetCurrentView(videoView); err != nil {
			return err
		}
	}

	return nil
}

func playVideo(g *gocui.Gui, v *gocui.View) error {
	cmd := exec.Command("bash", "-c",
		"mpv --ytdl-format='bestvideo[height<="+selectedVideoQuality+"]+bestaudio/best[height<="+selectedVideoQuality+"]' https://www.youtube.com/watch?v="+selectedVideo.ContentDetails.VideoId)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func downloadVideo(g *gocui.Gui, v *gocui.View) error {
	usr, _ := user.Current()
	outPath := filepath.FromSlash(usr.HomeDir + "/Downloads/")
	if _, err := os.Stat(config.Conf.DownloadPath); err == nil {
		if strings.HasSuffix(config.Conf.DownloadPath, "/") {
			outPath = config.Conf.DownloadPath
		} else {
			outPath = config.Conf.DownloadPath + "/"
		}
	}

	cmd := exec.Command("bash", "-c",
		"TS_SOCKET=/tmp/y-dl tsp youtube-dl -o '"+outPath+"%(title)s.%(ext)s' -f 'bestvideo[height<="+selectedVideoQuality+"]+bestaudio/best[height<="+selectedVideoQuality+"]' https://www.youtube.com/watch?v="+selectedVideo.ContentDetails.VideoId)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func selectQuality(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(qualityView, maxX/2-15, maxY/2-3, maxX/2+15, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Title = "Select Video quality"

		for _, quality := range videoQuality {
			fmt.Fprintln(v, quality)
		}

		if _, err := g.SetCurrentView(qualityView); err != nil {
			return err
		}
	}

	return nil
}

func rateVideo(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(rateVideoView, maxX/2-15, maxY/2-3, maxX/2+15, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Title = "Rate Video"

		for _, rating := range ratings {
			fmt.Fprintln(v, rating)
		}

		if _, err := g.SetCurrentView(rateVideoView); err != nil {
			return err
		}
	}

	return nil
}

func rate(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	e := api.RateVideo(selectedVideo.ContentDetails.VideoId, l)
	if e == nil {
		deleteRateVideo(g, v)
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
	if err := g.DeleteView(qualityView); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(videoView); err != nil {
		return err
	}

	return nil
}

func deleteRateVideo(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(rateVideoView); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(videoView); err != nil {
		return err
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
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(Layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
