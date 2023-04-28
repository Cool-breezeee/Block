package db

import (
	"github.com/boltdb/bolt"
	"log"
)

func BoltDbInit() *bolt.DB {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
