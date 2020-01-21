package main

import (
	"fmt"
	"github.com/kataras/iris"
)

func main()  {
	app := iris.New()

	app.Post("/upload", func(ctx iris.Context) {
		maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

		err := ctx.Request().ParseMultipartForm(maxSize)

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		form := ctx.Request().MultipartForm

		files := form.File["file"]
		//failures := 0

		for _, file := range files {
			fmt.Println(file.Filename)
		}
	})

	app.Run(iris.Addr(":8080"), iris.WithPostMaxMemory(32<<20))
}