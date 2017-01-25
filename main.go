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

	ServeAPI(ServerEnv{ db: db })
}
