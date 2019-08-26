package main

import (
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type User struct {
	Name string
}

func main()  {
	fmt.Println("Hello golang")

	app:= iris.New()
	mvc.Configure(app.Party("/doctor"), doctorMVC)

	app.Post("/decode", func(ctx iris.Context) {
		var user User
		// 请求参数格式化  请求参数是json类型转化成 User类型
		// 比如 post 参数 {username:'xxxx'} 转成 User 类型
		//把 json 类型请求参数 转成结构体
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s", user.Name)
	})

	app.Run(iris.Addr(":8090"))

}

func doctorMVC(app *mvc.Application)  {
	doctorService:= service.NewDoctorService()
	app.Register(doctorService)
	app.Handle(new(controllers.DoctorController))
}