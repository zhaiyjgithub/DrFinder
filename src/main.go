package main

import (
	"DrFinder/src/web/controllers"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main()  {
	fmt.Println("Hello golang")

	app:= iris.New()
	mvc.Configure(app.Party("/doctor"), doctorMVC)
	app.Run(iris.Addr(":8090"))

}

func doctorMVC(app *mvc.Application)  {
	app.Handle(new(controllers.DoctorController))
}