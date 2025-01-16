package db

import (
	"encoding/json"
	"rsscan/internal/common"

	"github.com/tidwall/buntdb"
)

// OpenDB opens a BuntDB database at the specified path.
func OpenDB(path string) (*buntdb.DB, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SaveFeed saves a feed with the specified key and value in the database.
func SaveFeed(db *buntdb.DB, key string, metadata common.PodcastMetadata) error {

	value, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(value), nil)
		return err
	})
}

// DeleteFeed deletes a feed with the specified key from the database.
func DeleteFeed(db *buntdb.DB, key string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
}

// GetFeed retrieves a feed with the specified key from the database.
func GetFeed(db *buntdb.DB, key string) (common.PodcastMetadata, error) {
	var metadata common.PodcastMetadata

	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		data, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(data), &metadata)
	})
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

// ListFeeds lists all the keys (feeds) in the database.
func ListFeeds(db *buntdb.DB) ([]common.PodcastMetadata, error) {
	var feeds []common.PodcastMetadata

	err := db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var feed common.PodcastMetadata
			err := json.Unmarshal([]byte(value), &feed)
			if err != nil {
				// TODO: print warning
				return true
			}
			feeds = append(feeds, feed)
			return true // continue iteration
		})
		return err
	})

	return feeds, err
}
