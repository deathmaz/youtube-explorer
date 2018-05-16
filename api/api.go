package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// Service service
var Service *youtube.Service

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// GetMySubscriptions get channels list for current user
func GetMySubscriptions() (res *youtube.SubscriptionListResponse, err error) {
	call := Service.Subscriptions.List("contentDetails,snippet")
	call = call.Mine(true)
	call = call.MaxResults(50)
	response, err := call.Do()
	handleError(err, "")

	return response, err
}

// GetChannelPlaylistItems get playlists for channel
func GetChannelPlaylistItems(channelID string) (res *youtube.PlaylistListResponse, err error) {
	call := Service.Playlists.List("contentDetails,snippet")
	call = call.ChannelId(channelID)
	call = call.MaxResults(50)
	response, err := call.Do()
	handleError(err, "")

	return response, err
}

// GetPlaylistItems get playlist items
func GetPlaylistItems(playlistID string) (res *youtube.PlaylistItemListResponse, err error) {
	call := Service.PlaylistItems.List("contentDetails,snippet")
	call = call.PlaylistId(playlistID)
	call = call.MaxResults(50)

	response, err := call.Do()
	handleError(err, "")
	return response, err
}

// GetChannel get channel by its id
func GetChannel(channelID string) (res *youtube.ChannelListResponse, err error) {
	call := Service.Channels.List("contentDetails,snippet")
	call = call.Id(channelID)
	response, err := call.Do()
	handleError(err, "")
	return response, err
}

// GetVideos get video by its id
func GetVideos(videoID string) (*youtube.VideoListResponse, error) {
	call := Service.Videos.List("statistics,contentDetails")
	call = call.Id(videoID)
	res, err := call.Do()
	handleError(err, "")

	return res, err
}

// RateVideo rate video
func RateVideo(videoID string, rating string) error {
	call := Service.Videos.Rate(videoID, rating)
	err := call.Do()
	return err
}

// GetYourRating get rating for video
func GetYourRating(videoID string) (string, error) {
	call := Service.Videos.GetRating(videoID)
	response, err := call.Do()
	handleError(err, "")
	return response.Items[0].Rating, err
}

// GetCommentThreads get comment threads
func GetCommentThreads(videoID string) (*youtube.CommentThreadListResponse, error) {
	call := Service.CommentThreads.List("snippet,replies")
	call = call.VideoId(videoID)
	call = call.MaxResults(50)
	call = call.TextFormat("plainText")
	response, err := call.Do()
	handleError(err, "")
	return response, err
}

// GetComment get comment for a thread
func GetComment(threadID string) (*youtube.Comment, error) {
	call := Service.Comments.List("snippet")

	call = call.Id(threadID)
	response, err := call.Do()
	handleError(err, "")

	return response.Items[0], err
}

// GetComments get comments for a video
func GetComments(videoID string) (*youtube.CommentListResponse, error) {
	threads, _ := GetCommentThreads(videoID)
	call := Service.Comments.List("snippet")

	comments := ""
	threadsLen := len(threads.Items) - 1
	for i, thread := range threads.Items {
		comments += thread.Id
		if threadsLen != i {
			comments += ","
		}
	}

	call = call.Id(comments)

	response, err := call.Do()
	handleError(err, "")

	return response, err
}

// GetReply get replies for comment
func GetReply(commentID string) (*youtube.CommentListResponse, error) {
	call := Service.Comments.List("snippet")
	call = call.ParentId(commentID)
	call = call.MaxResults(10)
	response, err := call.Do()
	handleError(err, "")
	return response, err
}

// Search search something
func Search(query string) {
	call := Service.Search.List("id,snippet").
		Q(query).
		MaxResults(50)
	response, err := call.Do()
	handleError(err, "")
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

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

// Run function
func init() {
	ctx := context.Background()
	usr, _ := user.Current()

	b, err := ioutil.ReadFile(filepath.FromSlash(usr.HomeDir + "/.config/my-youtube/client_secret.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")
	// channelsListByUsername(service, "contentDetails,snippet", "")

	Service = service
}
