package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rsscan/internal/common"
	"rsscan/internal/db"
	"rsscan/internal/rss"

	"github.com/tidwall/buntdb"
)

func addFeedCmd(database *buntdb.DB, rssURL string) error {
	metadata, err := rss.RequestRSSFeed(rssURL)
	if err != nil {
		return err
	}
	err = rss.AddRSSFeed(database, metadata)
	return err
}

func removeFeedCmd(database *buntdb.DB, rssURL string) error {
	err := rss.RemoveRSSFeed(database, rssURL)
	return err
}

func listFeedsCmd(database *buntdb.DB) error {
	feeds, err := rss.ListRSSFeeds(database)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n  Episode: %s\n  Published:%s\n",
			feed.ChannelTitle, feed.ItemTitle, feed.PubDate)
	}
	return nil
}

func updateEpisodesCmd(database *buntdb.DB) error {

	feeds, err := rss.ListRSSFeeds(database)
	if err != nil {
		log.Fatal(err)
	}

	// Refresh feed
	for _, feed := range feeds {
		latest, err := rss.RequestRSSFeed(feed.RSSURL)
		if err != nil {
			log.Fatal(err)
		}

		filePath := rss.BuildEpisodeName(feed.ChannelTitle)

		// TODO: do a proper time check
		if latest.PubDate != feed.PubDate {
			// refresh feed
			err := addFeedCmd(database, feed.RSSURL)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Updating.... %s\n  Episode: %s\n  Published:%s\n",
				feed.ChannelTitle, feed.ItemTitle, feed.PubDate)
			err = os.Remove(filePath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Check if the file exists, download if not
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			err := rss.DownloadLatestPodcast(&feed)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func getDBPath() string {
	path, err := common.GetDataDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/podcasts.db", path)
}

func printInfo() {
	path, err := common.GetDataDir()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=====================================")
	fmt.Println("  RSScan - A simple RSS feed reader")
	fmt.Printf("  Version:  %s\n", version)
	fmt.Printf("  Data dir: %s\n", path)
	fmt.Println("======================================")
	fmt.Println()
}

var version string

func main() {
	dbFile := getDBPath()
	database, err := db.OpenDB(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	addFeed := flag.String("add", "", "Add a new RSS feed")
	removeFeed := flag.String("remove", "", "Remove an RSS feed")
	listFeeds := flag.Bool("list", false, "List saved RSS feeds")
	info := flag.Bool("v", false, "Verbose")
	updateEpisodes := flag.Bool("update", false, "Update feed information and latest episode if needed")

	flag.Parse()

	if *info {
		printInfo()
	}

	if *addFeed != "" {
		err = addFeedCmd(database, *addFeed)
	} else if *removeFeed != "" {
		err = removeFeedCmd(database, *removeFeed)
	} else if *listFeeds {
		err = listFeedsCmd(database)
	} else if *updateEpisodes {
		err = updateEpisodesCmd(database)
	} else {
		flag.Usage()
	}

	if err != nil {
		log.Fatal(err)
	}
}
