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

	doctorParty := app.Party(Utils.APIDoctor, j.Serve)
	mvc.Configure(doctorParty, doctorMVC)

	userParty := app.Party(Utils.APIUser)
	mvc.Configure(userParty, userMVC)

	registerParty := app.Party(Utils.APIRegister)
	mvc.Configure(registerParty, registerMVC)

	app.Get("/SignIn", signIn)
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

func signIn(ctx iris.Context) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(conf.JWRTSecret)

	ctx.HTML(`Token: ` + tokenString + `<br/><br/>
<a href="/secured?token=` + tokenString + `">/secured?token=` + tokenString + `</a>`)
}
