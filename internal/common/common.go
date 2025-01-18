package common

import (
	"os"
)

type PodcastMetadata struct {
	ChannelTitle string
	ItemTitle    string
	RSSURL       string
	PubDate      string
	AudioURL     string
}

func GetDataDir() (string, error) {
	dataDir := os.Getenv("RSSCAN_DATA_DIR")
	if dataDir == "" {
		var err error
		dataDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return dataDir, nil
}
