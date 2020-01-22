package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PostController struct {
	Ctx iris.Context
	Service service.PostService
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.CreatePost, "CreatePost")
	b.Handle(iris.MethodPost, Utils.UpdatePost, "UpdatePost")
	b.Handle(iris.MethodPost,Utils.AddLikes, "AddLikes")
	b.Handle(iris.MethodPost, Utils.AddFavor, "AddFavor")
	b.Handle(iris.MethodPost, Utils.DeletePost, "DeletePost")
	b.Handle(iris.MethodPost, Utils.GetPostByPage, "GetPostByPage")
	b.Handle(iris.MethodGet, Utils.ImgPost, "ImgPost")
}

func (c *PostController) CreatePost() {
	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	req := c.Ctx.Request()
	err := req.ParseMultipartForm(maxSize)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "Max size", nil)
		return
	}

	userId, _ := strconv.Atoi(req.FormValue("UserID"))
	postType, _ := strconv.Atoi(req.FormValue("Type"))
	title := req.FormValue("Title")
	desc := req.FormValue("Description")

	var post models.Post
	post.Type = postType
	post.UserID = userId
	post.Title = title
	post.Description = desc
	post.Favorites = 0
	post.Likes = 0
	post.Priority = 0
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	err = c.Service.Add(&post)

	form := req.MultipartForm
	files := form.File["file"]
	failures := 0

	for _, file := range files {
		newFileName := fmt.Sprintf("%s-%s", generateFileName(5), file.Filename)
		_, err = saveFile(file, "./src/upload/", newFileName)

		if err != nil {
			failures = failures + 1
		}
	}

	if failures == len(files) {
		response.Fail(c.Ctx, response.Err, "save files fail", nil)
		return
	}

	if err != nil {
		response.Fail(c.Ctx, response.Err, "create post fail", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) UpdatePost()  {
	type Param struct {
		ID int `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var post models.Post
	post.ID = param.ID
	post.Description = param.Description

	err = c.Service.Update(&post)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "update post faiel", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddLikes()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddLike(param.ID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddFavor()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddFavor(param.ID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) DeletePost()  {
	type Param struct{
		ID int `validate:"gt=0"`
		UserID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.DeleteByUser(param.ID, param.UserID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) GetPostByPage()  {
	type Param struct {
		Type int
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	posts := c.Service.GetPostListByPage(param.Type, param.Page, param.PageSize)

	response.Success(c.Ctx, response.Successful, posts)
}

func (c *PostController) ImgPost()  {
	fileName := c.Ctx.URLParam("name")

	filePath := fmt.Sprintf("./src/upload/" + fileName)
	_ = c.Ctx.ServeFile(filePath, true)
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