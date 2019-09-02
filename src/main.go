package main

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/web/controllers"
	"crypto/md5"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

	app.RegisterView(iris.HTML("./src/web/templates/", ".html"))

	doctorParty := app.Party(Utils.APIDoctor, j.Serve)
	mvc.Configure(doctorParty, doctorMVC)

	userParty := app.Party(Utils.APIUser)
	mvc.Configure(userParty, userMVC)

	registerParty := app.Party(Utils.APIRegister)
	mvc.Configure(registerParty, registerMVC)

	app.Get("/upload", func(ctx iris.Context) {
		now := time.Now().Unix()
		h := md5.New()

		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		ctx.View("upload_form.html", token)
	})

	app.Post("/upload_manual", func(ctx 	iris.Context) {
		maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

		err := ctx.Request().ParseMultipartForm(maxSize)

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}

		form := ctx.Request().MultipartForm

		files := form.File["uploadfile"]
		failure := 0
		for _, file := range files {
			_, err = saveUploadedFiles(file, "./src/web/sources")
			if err != nil {
				failure ++
				ctx.Writef("failed to upload: %s\n", file.Filename)
			}
		}

		ctx.Writef("%d files uploaded", len(files) - failure)
	})


	app.Run(iris.Addr(":8090"))
}

func saveUploadedFiles(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(777))

	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, src)

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