package dataSource

import (
	"DrFinder/src/conf"
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"sync"
	"time"
)

const codeBucket = "CodeBucket"

var (
	masterEngine *gorm.DB
	cacheDB *bolt.DB
	mongoEngine *mongo.Database
	elasticSearchClient *elastic.Client

	masterDBOnce sync.Once
	mongoDBOnce sync.Once
	elasticSearchOnce sync.Once
)

func InstanceMaster() *gorm.DB {
	masterDBOnce.Do(func() {
		var err error
		c := conf.MasterDBConf
		driveSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			c.User, c.Pwd, c.Host, c.Port, c.DBName)

		masterEngine, err = gorm.Open(conf.DriverName,
			driveSource)

		if err != nil {
			log.Fatal("dbhelper instance error")
		}

		masterEngine.LogMode(true)
	})

	return masterEngine
}

func InstanceMongoDB() *mongo.Database {
	mongoDBOnce.Do(func() {
		conn := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?", conf.MongoDBConf.User, conf.MongoDBConf.Pwd,
			conf.MongoDBConf.Host, conf.MongoDBConf.Port, conf.MongoDBConf.DBName)
		clientOptions := options.Client().ApplyURI(conn)

		ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connected to MongoDB!")
		mongoEngine = client.Database(conf.MongoDBConf.DBName)
	})

	return mongoEngine
}

func InstanceElasticSearch()  *elastic.Client {
	elasticSearchOnce.Do(func() {
		//addr := net.JoinHostPort(conf.ElasticSearchConf.Host, conf.ElasticSearchConf.Port)
		url := fmt.Sprintf("%s/%d", conf.ElasticSearchConf.Host, conf.ElasticSearchConf.Port)
		client, err := elastic.NewClient(elastic.SetURL(url))
		if err != nil {
			log.Fatalf("Setup elastic search server failed....: %v", err)
		}

		info, code, err := client.Ping(url).Do(context.Background())
		if err != nil {
			log.Fatalf("ping elastic search server error: %v", err)
		}

		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		version, err := client.ElasticsearchVersion(url)
		if err != nil {
			log.Fatalf("Get elastic version error: %v", err)
		}
		fmt.Printf("Elasticsearch version %s\n", version)

		elasticSearchClient = client
	})

	return elasticSearchClient
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
