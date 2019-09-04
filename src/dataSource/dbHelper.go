package dataSource

import (
	"DrFinder/src/conf"
	"fmt"
	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	"log"
	"os"
	"sync"
)

const codeBucket = "CodeBucket"

var (
	masterEngine *gorm.DB
	cacheDB *bolt.DB
	lock sync.Mutex
)

func InstanceMaster() *gorm.DB {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	c := conf.MasterDBConf
	driveSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Pwd, c.Host, c.Port, c.DBName)

	fmt.Println(driveSource)
	engine, err := gorm.Open(conf.DriverName,
		driveSource)

	if err != nil {
		log.Fatal("dbhelper instance error")
	}

	engine.LogMode(true)

	return engine
}

func InstanceCacheDB() error {
	db, err := bolt.Open("./src/dataSource/bolt.db", os.FileMode(0750), nil)
	cacheDB = db

	if err != nil {
		return err
	}

	err = cacheDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(codeBucket))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func Save(key string, value string) error {
	err := cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(codeBucket))
		err := b.Put([]byte(key), []byte(value))

		return err
	})

	return err
}

func Get(key string) []byte {
	var v []byte

	 err := cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(codeBucket))
		v = b.Get([]byte(key))

		if v != nil {
			return nil
		}else {
			return errors.New("not found")
		}
	})

	if err != nil {
		return nil
	}

	return v
}
