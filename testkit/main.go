package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/boltdb/bolt/hack"
	"github.com/boltdb/bolt/log"
	"strconv"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Bolt.Error(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket(hack.Slice("myBucket"))
		bk := tx.Bucket(hack.Slice("myBucket"))
		return bk.Put(hack.Slice("key1"), hack.Slice("value1"))
	})

	db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(hack.Slice("myBucket"))
		v := bk.Get(hack.Slice("key1"))
		log.Bolt.Info(hack.String(v))
		return nil
	})
	
	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte("myBucket"))
	})

	db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucket(hack.Slice("myBucket"))
		if err != nil {
			return err
		}
		v := bk.Get([]byte("key1"))
		log.Bolt.Info(hack.String(v))
		return bk.Put(hack.Slice("key1"), hack.Slice("value1"))
	})
}

// createUser creates a new user in the given account.
func createUser(db *bolt.DB, accountID int, u *User) error {
	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Retrieve the root bucket for the account.
	// Assume this has already been created when the account was set up.
	root := tx.Bucket([]byte(strconv.FormatUint(uint64(accountID), 10)))

	// Setup the users bucket.
	bkt, err := root.CreateBucketIfNotExists([]byte("USERS"))
	if err != nil {
		return err
	}

	// Generate an ID for the new user.
	userID, err := bkt.NextSequence()
	if err != nil {
		return err
	}
	u.ID = userID

	// Marshal and save the encoded user.
	if buf, err := json.Marshal(u); err != nil {
		return err
	} else if err := bkt.Put([]byte(strconv.FormatUint(u.ID, 10)), buf); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
