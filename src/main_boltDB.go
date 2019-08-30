package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

const codeBucket = "codeBucket"

type Code struct {
	date time.Time
	value string
}

func main()  {
	db, err := bolt.Open("./src/dataSource/bolt.db", os.FileMode(0750), nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()


	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(codeBucket))
		if err != nil {
			return fmt.Errorf("create bucket error: %s", err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	k := "answer"
	v := "Zack"

	save(db, k, v)
	value := get(db, k)

	fmt.Printf("value: %s", value)
}

func save(db *bolt.DB, key string, value string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(codeBucket))
		err := b.Put([]byte(key), []byte(value))

		return err
	})

	return err
}

func get(db *bolt.DB, key string) string {
	var value string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(codeBucket))
		v := b.Get([]byte(key))

		value = string(v)
		fmt.Println(v)

		return nil
	})

	return value
}

