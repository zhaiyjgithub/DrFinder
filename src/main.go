package main

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"context"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	err := dataSource.InstanceCacheDB()

	if err != nil {
		panic(err)
	}

	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return conf.JWRTSecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, e error) {
			response.Fail(ctx, response.Expire, e.Error(), nil)
		},
	})

	app:= iris.New()

	//app.RegisterView(iris.HTML("./src/web/templates/", ".html"))

	doctorParty := app.Party(Utils.APIDoctor, j.Serve)
	mvc.Configure(doctorParty, doctorMVC)

	userParty := app.Party(Utils.APIUser)
	mvc.Configure(userParty, userMVC)

	registerParty := app.Party(Utils.APIRegister)
	mvc.Configure(registerParty, registerMVC)

	adminParty := app.Party(Utils.APIUtils)
	mvc.Configure(adminParty, adminMVC)

	advertisementParty := app.Party(Utils.APIAd)
	mvc.Configure(advertisementParty, advertiseMVC)

	postParty := app.Party(Utils.APIPost)
	mvc.Configure(postParty, postMVC)

	answerParty := app.Party(Utils.APIAnswer)
	mvc.Configure(answerParty, answerMVC)

	app.Run(iris.Addr(":8090"))


}

func doctorMVC(app *mvc.Application) {
	service := service.NewDoctorService()
	app.Register(service)
	app.Handle(new(controllers.DoctorController))
}

func userMVC(app *mvc.Application)  {
	service := service.NewUserService()
	app.Register(service)
	app.Handle(new(controllers.UserController))
}

func registerMVC(app *mvc.Application)  {
	service := service.NewUserService()
	app.Register(service)
	app.Handle(new(controllers.RegisterController))
}

func adminMVC(app *mvc.Application)  {
	app.Handle(new(controllers.AdminController))
}

func advertiseMVC(app *mvc.Application)  {
	service := service.NewAdvertiseService()
	app.Register(service)
	app.Handle(new(controllers.AdvertisementController))
}

func postMVC(app *mvc.Application)  {
	service := service.NewPostService()
	app.Register(service)
	app.Handle(new(controllers.PostController))
}

func answerMVC(app *mvc.Application)  {
	service := service.NewAnswerService()
	app.Register(service)
	app.Handle(new(controllers.AnswerController))
}

func testLog()  {
	log := logrus.New()

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size": 10,
	}).Info("A group of walrus emerges from the ocean")
}

func testMongoDB() {
	type Trainer struct {
		Name string
	}

	mongoURI := "mongodb://myTester:123456@localhost:27017/test?"

	// Set client options
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
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)

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

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	fmt.Println(cur)
}

