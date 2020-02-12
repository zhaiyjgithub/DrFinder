package controllers

import (
	"DrFinder/src/utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"io/ioutil"
	"path/filepath"
	"time"
)

type PostController struct {
	Ctx iris.Context
	Service service.PostService
	AnswerService service.AnswerService
	UserService service.UserService
	PostImageService service.PostImageService
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, utils.CreatePost, "CreatePost")
	b.Handle(iris.MethodPost, utils.UpdatePost, "UpdatePost")
	b.Handle(iris.MethodPost, utils.AddLikes, "AddLikes")
	b.Handle(iris.MethodPost, utils.AddFavor, "AddFavor")
	b.Handle(iris.MethodPost, utils.DeletePost, "DeletePost")
	b.Handle(iris.MethodPost, utils.GetPostByPage, "GetPostByPage")
	b.Handle(iris.MethodGet, utils.ImgPost, "ImgPost")
}

func (c *PostController) UpdatePost()  {
	type Param struct {
		ID int `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var post models.Post
	post.ID = param.ID
	post.Description = param.Description

	err = c.Service.Update(&post)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "update post fail", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddLikes()  {
	type Param struct {
		UserID int `validate:"gt=0"`
		PostID int `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddLike(param.PostID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddFavor()  {
	type Param struct {
		UserID int `validate:"gt=0"`
		PostID int `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddFavor(param.PostID)

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

	err := utils.ValidateParam(c.Ctx, validate, &param)

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

	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	type PostInfo struct {
		PostID int
		UserIcon string
		UserName string
		Type int
		Title string
		Description string
		Likes int
		PostDate time.Time
		AnswerCount int
		LastAnswerName string
		LastAnswerDate time.Time
		URLs []string
	}

	posts := c.Service.GetPostListByPage(param.Type, param.Page, param.PageSize)

	var postInfos []PostInfo
	for i := 0; i < len(posts); i ++  {
		post := &posts[i]
		postUser := c.UserService.GetUserById(post.UserID)
		answer, count := c.AnswerService.GetLastAnswer(post.ID)

		var answerName string
		var lastCreateAt time.Time

		if answer != nil {
			user := c.UserService.GetUserById(answer.UserID)
			answerName = user.FirstName
			lastCreateAt = answer.CreatedAt
		}

		postImgs := c.PostImageService.GetImageByPostId(post.ID)

		urls := make([]string, 0)
		for _, img:= range postImgs {
			urls = append(urls, img.URL)
		}

		var postInfo PostInfo
		postInfo.PostID = post.ID
		postInfo.UserIcon = postUser.HeaderIcon
		postInfo.UserName = postUser.LastName
		postInfo.Type = post.Type
		postInfo.Title = post.Title
		postInfo.Description = post.Description
		postInfo.Likes = post.Likes
		postInfo.PostDate = post.CreatedAt
		postInfo.AnswerCount = count
		postInfo.LastAnswerName = answerName
		postInfo.LastAnswerDate = lastCreateAt
		postInfo.URLs = urls
		postInfos = append(postInfos, postInfo)
	}

	response.Success(c.Ctx, response.Successful, postInfos)
}

func (c *PostController) ImgPost()  {
	fileName := c.Ctx.URLParam("name")

	filePath := fmt.Sprintf("./src/upload/post/" + fileName)
	_ = c.Ctx.ServeFile(filePath, true)
}

func (c *PostController) CreatePost()  {
	type Param struct { //`validate:"gt=0"`
		UserID int `validate:"gt=0"`
		Type int `validate:"numeric"`
		Title string `validate:"gt=0"`
		Description string `validate:"gt=0"`
		Files []struct{
			Ext string `validate:"gt=0"`
			Base64Data string `validate:"base64"`
		}
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var post models.Post
	post.Type = param.Type
	post.UserID = param.UserID
	post.Title = param.Title
	post.Description = param.Description
	post.Favorites = 0
	post.Likes = 0
	post.Priority = 0
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	err, postId := c.Service.Add(&post)

	files := param.Files
	failure := 0
	for _, file := range files {
		fileName := fmt.Sprintf("%s.%s", generateFileName(param.UserID), file.Ext)
		_, err = saveFile(file.Base64Data, "./src/upload/post", fileName)
		if err != nil {
			failure ++
		}else {
			postImg := models.PostImage{PostID:postId, URL: fileName}
			_ = c.PostImageService.CreatePostImage(postImg)
		}
	}

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func saveFile(imgBase64 string, destDir string, fileName string) (int64, error)  {
	imgBuffer, err  := base64.StdEncoding.DecodeString(imgBase64)

	if err != nil {
		fmt.Println("decode img base64 error")

		return 1, err
	}

	filePath := filepath.Join(destDir, fileName)
	err = ioutil.WriteFile(filePath, imgBuffer, 0666)

	if err != nil {
		return 1, err
	}

	return 0, nil
}

func generateFileName(userId int) string {
	data := []byte(fmt.Sprintf("%d-%d", userId, time.Now().Unix()))
	md5er := md5.New()
	md5er.Write(data)

	return hex.EncodeToString(md5er.Sum(nil))
}