package rss

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"rsscan/internal/common"
	"rsscan/internal/db"

	"github.com/mmcdole/gofeed"
	"github.com/tidwall/buntdb"
)

func RequestRSSFeed(rssURL string) (*common.PodcastMetadata, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		return nil, err
	}

	// Validate that items are prsent
	if len(feed.Items) <= 0 {
		return nil, errors.New("no items present in feed")
	}

	metadata := common.PodcastMetadata{
		ChannelTitle: feed.Title,
		ItemTitle:    feed.Items[0].Title,
		RSSURL:       rssURL,
		PubDate:      feed.Items[0].Published,
		AudioURL:     feed.Items[0].Enclosures[0].URL,
	}

	return &metadata, nil
}

func AddRSSFeed(database *buntdb.DB, metadata *common.PodcastMetadata) error {
	return db.SaveFeed(database, metadata.RSSURL, *metadata)
}

func RemoveRSSFeed(database *buntdb.DB, rssURL string) error {
	return db.DeleteFeed(database, rssURL)
}

func ListRSSFeeds(database *buntdb.DB) ([]common.PodcastMetadata, error) {
	return db.ListFeeds(database)
}

// generate quick hash for episode name and limit to 32 characters
func BuildEpisodeName(title string) string {
	outName := ""
	for _, char := range title {
		outName += fmt.Sprintf("%02X", char)
	}
	return outName[:32]
}

func DownloadLatestPodcast(feedData *common.PodcastMetadata) error {

	// Create the episodes directory if it doesn't exist
	if err := os.MkdirAll("episodes", os.ModePerm); err != nil {
		return err
	}

	// Download the latest episode
	response, err := http.Get(feedData.AudioURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create a file to save the episode
	filePath := BuildEpisodeName(feedData.ChannelTitle)
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the response body to the file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Donwloading: %s\n   Episode: %s\n", feedData.ChannelTitle, feedData.ItemTitle)
	return nil
}
