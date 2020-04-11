package main

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main()  {
	//testMongoDB()

	//findOne()

	type Event struct {
		Name string
		Time time.Time
	}

	dao := dao.NewMongoDao(dataSource.InstanceMongoDB())
	//event := Event{Name: "Login", Time: time.Now()}
	//
	//err := dao.InsertOne("event", event)
	//if err != nil {
	//	fmt.Println("insert failed")
	//}else {
	//	fmt.Println("insert success")
	//}

	res, err := dao.FindOne("event", bson.M{"name": "Login"})
	fmt.Println(res)
	fmt.Println(err)
}

func testMongoDB() {
	type Trainer struct {
		//ID   primitive.ObjectID `bson:"_id"`
		Name string `bson:"name"`
	}

	//mongoURI := "mongodb://myTester:123456@localhost:27017"

	//credential := options.Credential{
	//	Username: "myUserAdmin",
	//	Password: "123456",
	//	AuthSource: "admin",
	//}

	// Set client options

	//mongo userTrack -u myUserAdmin -p 123456 --authenticationDatabase admin
	//setAuth 就是使用admin user  to login other database.

	mongoURI := "mongodb://myTester:123456@localhost:27017/test?"
	clientOptions := options.Client().ApplyURI(mongoURI)


	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("mycol")

	findOptions := options.Find()
	findOptions.SetLimit(1)

	//var t Trainer
	//t.Name = "Mary"

	cur, err := collection.Find(context.TODO(), findOptions)

	if err != nil {
		log.Fatal(err)
	}

	var results []Trainer
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	//p := results[0]
	//
	//fmt.Printf("Found multiple documents (array of pointers): %+v\n",p.Name)

	fmt.Println(cur)
}

func addOne()  {

	type Trainer struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string `bson:"name"`
	}

	col := dataSource.InstanceMongoDB().Collection("mycol")

	var t Trainer
	t.Name = "Mary col123"

	res, err := col.InsertOne(context.TODO(), t)

	fmt.Println(res)
	fmt.Println(err)
}

func findOne() error {
	b, err := bson.Marshal(bson.D{{"foo", "bar"}})
	if err != nil { return err }
	var fooer struct {
		Foo string
	}
	err = bson.Unmarshal(b, &fooer)
	if err != nil { return err }

	return nil
}


