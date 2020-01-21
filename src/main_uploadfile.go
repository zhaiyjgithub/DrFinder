package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kataras/iris"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func main()  {
	app := iris.New()

	app.Get("/img", func(ctx iris.Context) {
		//fileContentType := "image/jpeg"
		fileName := ctx.URLParam("name")
		//ctx.ContentType(fileContentType)

		filePath := fmt.Sprintf("./src/upload/" + fileName)
		ctx.ServeFile(filePath, true)
	})

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
		failures := 0

		for _, file := range files {
			formatFileName := fmt.Sprintf("%s-%s", generateFileName(5), file.Filename)
			_, err := saveFile(file, "./src/upload", formatFileName)

			if err != nil {
				failures = failures + 1
			}
		}

		fmt.Println(failures)
	})

	app.Run(iris.Addr(":8080"), iris.WithPostMaxMemory(32<<20))
}

func saveFile(fh *multipart.FileHeader, destDir string, fileName string) (int64, error)  {
	src, err := fh.Open()

	if err != nil {
		return 0, err
	}

	defer src.Close()

	out, err := os.OpenFile(filepath.Join(destDir, fileName), os.O_WRONLY | os.O_CREATE , os.FileMode(0666))

	if err != nil {
		return 0, err
	}

	defer out.Close()

	return io.Copy(out, src)
}

func generateFileName(userId int) string {
	data := []byte(fmt.Sprintf("%d-%d", userId, time.Now().Unix()))
	md5er := md5.New()
	md5er.Write(data)
	
	return hex.EncodeToString(md5er.Sum(nil))
}