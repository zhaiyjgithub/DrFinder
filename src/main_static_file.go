package main

import (
	"github.com/kataras/iris"
)

func main()  {
	app := iris.New()

	//app.Get("/img", func(ctx iris.Context) {
	//	fileContentType := "image/jpeg"
	//
	//	ctx.ContentType(fileContentType)
	//	ctx.ServeFile("./src/web/sources/test.png", true)
	//})

	app.StaticWeb("/static", "./src/web/sources/")

	app.Run(iris.Addr(":8090"))
}


