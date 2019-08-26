package main

import (
	"DrFinder/src/Utils"
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main()  {
	fmt.Println("Hello golang")

	app:= iris.New()
	mvc.Configure(app.Party(Utils.APIDoctor), doctorMVC)

	app.Run(iris.Addr(":8090"))

}

func doctorMVC(app *mvc.Application)  {
	doctorService:= service.NewDoctorService()
	app.Register(doctorService)
	app.Handle(new(controllers.DoctorController))
}