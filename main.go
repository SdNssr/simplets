package main

import (
	"github.com/boltdb/bolt"
	"time"
	"log"
	"fmt"
)

func main() {
	db, err := bolt.Open("ts.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("timeseries"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	ServeAPI(ServerEnv{ db: db })
}
