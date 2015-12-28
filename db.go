package main

import (
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

type database struct {
	DB *bolt.DB
}

func initDB(filename string) database {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return database{DB: db}
}

func (db *database) Close() {
	db.DB.Close()
}

func (db *database) get(pubkey []byte) error {
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		if b == nil {
			return errors.New("Could not find bucket")
		}

		v := b.Get(pubkey)
		if v == nil {
			return errors.New("Could not find pubkey")
		}

		return nil
	})
	return err
}

func (db *database) set(pubkey []byte) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return b.Put(pubkey, []byte{0x01})
	})

	return err
}
