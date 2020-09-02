package controllers

import (
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/utils"
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
	AppendService service.AppendService
	ElasticService service.PostElasticService
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, utils.CreatePost, "CreatePost")
	b.Handle(iris.MethodPost, utils.UpdatePost, "UpdatePost", j.Serve)
	b.Handle(iris.MethodPost, utils.AddLikes, "AddLikes", j.Serve)
	b.Handle(iris.MethodPost, utils.AddFavor, "AddFavor", j.Serve)
	b.Handle(iris.MethodPost, utils.DeletePost, "DeletePost", j.Serve)
	b.Handle(iris.MethodPost, utils.GetPostByPage, "GetPostByPage")
	b.Handle(iris.MethodGet, utils.ImgPost, "ImgPost")
	b.Handle(iris.MethodPost, utils.GetMyPostByPage, "GetMyPostByPage", j.Serve)
	b.Handle(iris.MethodPost, utils.AddAppendToPost, "AddAppendToPost", j.Serve)
	b.Handle(iris.MethodPost, utils.GetAppendByPostID, "GetAppendByPostID", j.Serve)
	b.Handle(iris.MethodPost, utils.DeletePostByIds, "DeletePostByIds", j.Serve)
	b.Handle(iris.MethodPost, utils.SearchPostByPageFromElastic, "SearchPostByPageFromElastic")
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

func (c *PostController) DeletePostByIds()  {
	type Param struct {
		IDs []int `validate:"gt=0"`
		UserID int `validate:"gt=0"`
	}

	var p Param
	err := utils.ValidateParam(c.Ctx, validate, &p)
	if err != nil {
		return
	}

	for _, ID := range p.IDs {
		_ = c.Service.DeleteByUser(ID, p.UserID)
	}

	response.Success(c.Ctx, response.Successful, nil)
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

	//TO-DO 使用联结查询优化 view-model
	posts := c.Service.GetPostListByPage(param.Type, param.Page, param.PageSize)

	var postInfos []PostInfo
	for i := 0; i < len(posts); i ++  {
		post := posts[i]
		postUser := c.UserService.GetUserById(post.ID)
		answer, count := c.AnswerService.GetLastAnswer(post.ID)

		var answerName string
		var lastCreateAt time.Time

		if answer != nil {
			user := c.UserService.GetUserById(answer.UserID)
			answerName = user.Name
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
		postInfo.UserName = postUser.Name
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

func (c *PostController) GetMyPostByPage()  {
	type Param struct {
		UserID int `validate:"gt=0"`
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

	//TO-DO 使用联结查询优化 view-model
	posts := c.Service.GetMyPostListByPage(param.UserID, param.Page, param.PageSize)

	var postInfos []PostInfo
	for i := 0; i < len(posts); i ++  {
		post := posts[i]
		postUser := c.UserService.GetUserById(post.UserID)
		answer, count := c.AnswerService.GetLastAnswer(post.ID)

		var answerName string
		var lastCreateAt time.Time
		if answer != nil {
			user := c.UserService.GetUserById(answer.UserID)
			answerName = user.Name
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
		postInfo.UserName = postUser.Name
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

func (c *PostController) SearchPostByPageFromElastic()  {
	type Param struct {
		Content string `validate:"gt=0"`
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

	postIds := c.ElasticService.QueryPost(param.Content, param.Page, param.PageSize)
	posts := c.Service.GetPostByPostId(postIds)

	var postInfos []PostInfo
	for i := 0; i < len(posts); i ++  {
		post := posts[i]
		postUser := c.UserService.GetUserById(post.ID)
		answer, count := c.AnswerService.GetLastAnswer(post.ID)

		var answerName string
		var lastCreateAt time.Time

		if answer != nil {
			user := c.UserService.GetUserById(answer.UserID)
			answerName = user.Name
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
		postInfo.UserName = postUser.Name
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

func (c *PostController) AddAppendToPost()  {
	type Param struct {
		PostID int `validate:"gt=0"`
		Append string `validate:"gt=0"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	append := &models.Append{PostID: param.PostID, Content:param.Append}

	err = c.AppendService.AddAppend(append)
	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, "", nil)
	}
}

func (c *PostController) GetAppendByPostID()  {
	type Param struct {
		PostID int `validate:"gt=0"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	appends := c.AppendService.GetAppends(param.PostID)

	response.Success(c.Ctx, response.Successful, appends)
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
	post.ID = postId

	//files := param.Files
	//failure := 0
	//for _, file := range files {
	//	fileName := fmt.Sprintf("%s.%s", generateFileName(param.UserID), file.Ext)
	//	_, err = saveFile(file.Base64Data, "./src/upload/post", fileName)
	//	if err != nil {
	//		failure ++
	//	}else {
	//		postImg := models.PostImage{PostID:postId, URL: fileName}
	//		_ = c.PostImageService.CreatePostImage(postImg)
	//	}
	//}

	//sync post to elastic
	err = c.syncPostToElastic(&post)
	if err != nil {
		fmt.Printf("Sync post to elastic failed. Id: %d", post.ID)
	}

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) syncPostToElastic(post *models.Post) error {
	return c.ElasticService.AddOnePost(post)
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