package main

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
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

	app.Get("/SignIn", signIn)
	app.Run(iris.Addr(":8090"))
}

func doctorMVC(app *mvc.Application) {
	doctorService := service.NewDoctorService()
	app.Register(doctorService)
	app.Handle(new(controllers.DoctorController))
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
