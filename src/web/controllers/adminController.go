package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/response"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type AdminController struct {
	Ctx iris.Context
}

var adminValidate *validator.Validate

func (c *AdminController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.UploadFile, "UploadFile")
}

func (c *AdminController) UploadFile()  {
	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err := c.Ctx.Request().ParseMultipartForm(maxSize)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "parse file failed", nil)
		return
	}

	form := c.Ctx.Request().MultipartForm

	files := form.File["file"]
	failure := 0
	for _, file := range files {
		_, err = saveUploadedFiles(file, "./src/web/sources")
		if err != nil {
			failure ++
		}
	}

	if len(files) - failure > 0 {
		response.Success(c.Ctx, response.Successful, nil)
	}else {
		response.Fail(c.Ctx, response.Err, "parse file failed", nil)
	}
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