package main

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/sirupsen/logrus"
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
