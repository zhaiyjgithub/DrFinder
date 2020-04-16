package controllers

import (
	"DrFinder/src/conf"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/utils"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type UserController struct {
	Ctx iris.Context
	UserService service.UserService
	DoctorService service.DoctorService
	CollectionService service.CollectionService
	AnswerService service.AnswerService
	PostImageService service.PostImageService
	FeedbackService service.FeedbackService
}

var userValidate *validator.Validate

func (c *UserController) BeforeActivation(b mvc.BeforeActivation)  {
	userValidate = validator.New()

	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return conf.JWRTSecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, e error) {
			response.Fail(ctx, response.Expire, e.Error(), nil)
		},
	})

	b.Handle(iris.MethodPost, utils.CreateUser, "CreateUser")
	b.Handle(iris.MethodPost, utils.UpdatePassword, "UpdatePassword", j.Serve)
	b.Handle(iris.MethodPost, utils.UpdateUserInfo, "UpdateUserInfo", j.Serve)
	b.Handle(iris.MethodPost, utils.GetUserInfo, "GetUserInfo")
	b.Handle(iris.MethodPost, utils.GetMyFavorite, "GetMyFavorite", j.Serve)
	b.Handle(iris.MethodPost, utils.AddFavorite, "AddFavorite", j.Serve)
	b.Handle(iris.MethodPost, utils.DeleteMyFavorite, "DeleteMyFavorite", j.Serve)
	b.Handle(iris.MethodPost, utils.AddNewFeedback, "AddNewFeedback", j.Serve)
}

func (c *UserController) CreateUser() {
	type Param struct {
		Email      string `validate:"email"`
		Password   string `validate:"min=8,max=20"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	var user models.User
	user.Email = param.Email
	user.Password = param.Password

	err = c.UserService.CreateUser(&user)

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) UpdatePassword() {
	type Param struct {
		Email string `validate:"email"`
		OldPwd string `validate:"min=6,max=20"`
		NewPwd string `validate:"min=6,max=20"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	err = c.UserService.UpdatePassword(param.Email, param.OldPwd, param.NewPwd)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "email or old password is wrong", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) UpdateUserInfo()  {
	type Param struct {
		UserID int `validate:"gt=0"`
		Name string `validate:"gt=5"` //name length >= 6
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)
	if err != nil {
		return
	}

	var user models.User
	user.ID = param.UserID
	user.Name = param.Name

	err = c.UserService.UpdateUser(&user)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) GetUserInfo() {
	type Param struct {
		UserID int `validate:"gt=0"`
	}
	
	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}
	
	user := c.UserService.GetUserById(param.UserID)
	response.Success(c.Ctx, response.Successful, *user)
}

func (c *UserController) AddFavorite()  {
	type Param struct {
		UserID int
		ObjectID int
		ObjectType int
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	err = c.CollectionService.Add(param.UserID, param.ObjectID, param.ObjectType)
	if err != nil {
		errCode := response.Err
		if err.Error() == "is existing" {
			errCode = response.IsExist
		}

		response.Fail(c.Ctx, errCode, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) GetMyFavorite()  {
	type Param struct {
		UserID int
		Type int
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	if param.Type == 0 {
		doctors := c.CollectionService.GetMyFavoriteDoctors(param.UserID, param.Type, param.Page, param.PageSize)
		response.Success(c.Ctx, response.Successful, doctors)
	}else {
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

		posts := c.CollectionService.GetMyFavoritePosts(param.UserID, param.Type, param.Page, param.PageSize)

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
}

func (c *UserController) AddNewFeedback()  {
	type Param struct {
		UserID int `validate:"gt=0"`
		Feedback string `validate:"gt=0"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	var fb models.Feedback
	fb.UserID = param.UserID
	fb.Content = param.Feedback

	err = c.FeedbackService.AddFeedback(&fb)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController)DeleteMyFavorite()  {
	type Param struct {
		UserID int
		ObjectIDs []int `validate:'gt=0'`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	err = c.CollectionService.DeleteMyFavorite(param.UserID, param.ObjectIDs)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}
