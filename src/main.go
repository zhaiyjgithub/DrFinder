package main

import (
	"DrFinder/src/dataSource"
	"DrFinder/src/service"
	"DrFinder/src/utils"
	"DrFinder/src/web/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/sirupsen/logrus"
)

func main() {
	err := dataSource.InstanceCacheDB()

	if err != nil {
		panic(err)
	}
	//j := jwt.New(jwt.Config{
	//	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	//		return conf.JWRTSecret, nil
	//	},
	//	SigningMethod: jwt.SigningMethodHS256,
	//	ErrorHandler: func(ctx iris.Context, e error) {
	//		response.Fail(ctx, response.Expire, e.Error(), nil)
	//	},
	//})

	app:= iris.New()

	//app.RegisterView(iris.HTML("./src/web/templates/", ".html"))

	//doctorParty := app.Party(utils.APIDoctor, j.Serve)

	doctorParty := app.Party(utils.APIDoctor)
	mvc.Configure(doctorParty, doctorMVC)

	userParty := app.Party(utils.APIUser)
	mvc.Configure(userParty, userMVC)

	registerParty := app.Party(utils.APIRegister)
	mvc.Configure(registerParty, registerMVC)

	adminParty := app.Party(utils.APIUtils)
	mvc.Configure(adminParty, adminMVC)

	advertisementParty := app.Party(utils.APIAd)
	mvc.Configure(advertisementParty, advertiseMVC)

	postParty := app.Party(utils.APIPost)
	mvc.Configure(postParty, postMVC)

	answerParty := app.Party(utils.APIAnswer)
	mvc.Configure(answerParty, answerMVC)

	userTrackParty := app.Party(utils.APIUserTrack)
	mvc.Configure(userTrackParty, userTrackMVC)

	_ = app.Run(iris.Addr(":8090"), iris.WithPostMaxMemory(32<<20)) //max = 32M
}

func doctorMVC(app *mvc.Application) {
	doctorService := service.NewDoctorService()
	geoService := service.NewGeoService()
	affiliationService := service.NewAffiliationService()
	awardService := service.NewAwardService()
	cerService := service.NewCertificationService()
	clinicService := service.NewClinicalService()
	eduService := service.NewEducationService()
	langService := service.NewLangService()
	memberShipService := service.NewMembershipService()
	collectionService := service.NewCollectionService()

	app.Register(
		doctorService,
		affiliationService,
		awardService,
		cerService,
		clinicService,
		eduService,
		geoService,
		langService,
		memberShipService,
		collectionService,
		)
	app.Handle(new(controllers.DoctorController))
}

func userMVC(app *mvc.Application)  {
	userService := service.NewUserService()
	doctorService := service.NewDoctorService()
	collectionService := service.NewCollectionService()
	answerService := service.NewAnswerService()
	postImageService := service.NewPostImageService()
	feedbackService := service.NewFeedbackService()

	app.Register(
		userService,
		doctorService,
		collectionService,
		answerService,
		postImageService,
		feedbackService,
		)
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
	postService := service.NewPostService()
	answerService := service.NewAnswerService()
	userService := service.NewUserService()
	postImageService := service.NewPostImageService()
	appendService := service.NewAppendService()

	app.Register(
		postService,
		answerService,
		userService,
		postImageService,
		appendService,
		)
	app.Handle(new(controllers.PostController))
}

func answerMVC(app *mvc.Application)  {
	answerService := service.NewAnswerService()
	userService := service.NewUserService()
	app.Register(
			answerService,
			userService,
		)
	app.Handle(new(controllers.AnswerController))
}
func userTrackMVC(app *mvc.Application)  {
	userTrackService := service.NewUserTrackService()
	app.Register(
			userTrackService,
		)

	app.Handle(new(controllers.UserTrackController))
}

func testLog()  {
	log := logrus.New()

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size": 10,
	}).Info("A group of walrus emerges from the ocean")
}
