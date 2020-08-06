package controllers

import (
	"DrFinder/src/utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/o1egl/govatar"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

type RegisterController struct {
	Ctx iris.Context
	Service service.UserService
}

type VerificationCode struct {
	Email string `json:"email"`
	DateTime string `json:"date"`
	Value string `json:"value"`
}

var registerValidate *validator.Validate

func (c *RegisterController) BeforeActivation(b mvc.BeforeActivation)  {
	registerValidate = validator.New()
	b.Handle(iris.MethodPost, utils.SendVerificationCode, "SendVerificationCode")
	b.Handle(iris.MethodPost, utils.Register, "Register")
	b.Handle(iris.MethodPost, utils.SignIn, "SignIn")
	b.Handle(iris.MethodPost, utils.ResetPassword, "ResetPassword")
	b.Handle(iris.MethodPost, utils.VerifyEmail, "VerifyEmail")
	b.Handle(iris.MethodGet, utils.HeaderImg, "HeaderImg")
	b.Handle(iris.MethodGet, utils.DownloadFile, "DownloadFile")
}

func (c *RegisterController) DownloadFile()  {
	//
}

func (c *RegisterController) HeaderImg()  {
	fileName := c.Ctx.URLParam("name")

	filePath := fmt.Sprintf("./src/upload/user/" + fileName)
	_ = c.Ctx.ServeFile(filePath, true)
}

func (c *RegisterController) SendVerificationCode() {
	type Param struct {
		Email string `validate:"email"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, registerValidate, &param)
	if err != nil  {
		return
	}

	v := getCode()
	var vcode = &VerificationCode{
		Email: param.Email,
		DateTime:  time.Now().Format(conf.TimeFormat),
		Value: v,
	}

	cb, _ := json.Marshal(vcode)
	_ = dataSource.Save(param.Email, string(cb))

	err = sendEmail(param.Email, vcode.Value)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "send email fail", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *RegisterController) Register() {
	type Param struct {
		Email string `validate:"email"`
		Password string `validate:"min=6,max=20"`
		VerificationCode string `validate:"len=6"`
		Name string
		Gender int
	}
	
	var param Param
	err := utils.ValidateParam(c.Ctx, registerValidate, &param)
	if err != nil {
		return
	}

	_, err = c.Service.GetUserByEmail(param.Email)
	if err == nil {
		response.Fail(c.Ctx, response.IsExist, "This email has been registered", nil)
	}

	var code = VerificationCode{}
	cb := dataSource.Get(param.Email)
	if cb == nil {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
		return
	}

	err = json.Unmarshal(cb, &code)
	t, err := time.Parse(conf.TimeFormat, code.DateTime)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "parse time error", nil)
		return
	}

	tstamp := t.Unix()
	t2, err := time.Parse(conf.TimeFormat, time.Now().Format(conf.TimeFormat))
	nstamp := t2.Unix()

	if nstamp - tstamp > 60*10 {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
	}else if param.VerificationCode != code.Value {
		response.Fail(c.Ctx, response.Err, "verification code is wrong", nil)
	}else {
		var user models.User
		user.Email = param.Email
		user.Password = param.Password
		user.Name = param.Name
		//user.Gender = param.Gender

		fileName := utils.GenerateFileName(1)
		filePath := fmt.Sprintf("%s%s%s", "./src/upload/user/", fileName, ".jpg")
		err = govatar.GenerateFile(govatar.MALE, filePath)
		if err == nil {
			user.HeaderIcon = fmt.Sprintf("%s%s", fileName, ".jpg")
		}

		err = c.Service.CreateUser(&user)
		if err != nil {
			response.Fail(c.Ctx, response.Err, "create user failed", nil)
			return
		}

		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *RegisterController) SignIn()  {
	type Param struct {
		Email string `validate:"email"`
		Password string `validate:"gt=6"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, registerValidate, &param)
	if err != nil {
		return
	}

	user, err := c.Service.GetUserByEmail(param.Email)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "not exist", nil)
	}else if user.Password != param.Password {
		response.Fail(c.Ctx, response.Err, "password is wrong", nil)
	}else {
		token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"dispatch_time": time.Now().Format(conf.TimeFormat),
		})
		tokenString, _ := token.SignedString(conf.JWRTSecret)

		type UserInfo struct {
			*models.User
			Token string
		}

		var userInfo UserInfo
		userInfo.User = user
		userInfo.Token = tokenString
		response.Success(c.Ctx, "login success", userInfo)
	}
}

func (c * RegisterController) VerifyEmail()  {
	type Param struct {
		Email string `validate:"email"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	user, err := c.Service.GetUserByEmail(param.Email)
	if err == nil {
		response.Success(c.Ctx, "", user)
	}else {
		response.Success(c.Ctx, err.Error(), nil)
	}
}

func (c *RegisterController) ResetPassword()  {
	type Param struct {
		Email string `validate:"email"`
		Password string `validate:"min=6,max=20"`
		VerificationCode string `validate:"len=6"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	var code = VerificationCode{}
	cb := dataSource.Get(param.Email)
	if cb == nil {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
		return
	}

	err = json.Unmarshal(cb, &code)
	t, err := time.Parse(conf.TimeFormat, code.DateTime)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "parse time error", nil)
		return
	}

	tstamp := t.Unix()
	t2, err := time.Parse(conf.TimeFormat, time.Now().Format(conf.TimeFormat))
	nstamp := t2.Unix()

	if nstamp - tstamp > 60*10 {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
	}else if param.VerificationCode != code.Value {
		response.Fail(c.Ctx, response.Err, "verification code is wrong", nil)
	}else {
		err = c.Service.ResetPassword(param.Email, param.Password)
		if err != nil {
			response.Fail(c.Ctx, response.Err, "update failed", nil)
		}else {
			response.Success(c.Ctx, response.Successful, nil)
		}
	}
}

func sendEmail(email string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.ServerEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "DrFinder Verification code")

	body := fmt.Sprintf("Your verification code: %s", code)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(conf.Smtp, 587, conf.ServerEmail, conf.ServerEmailPwd)

	return d.DialAndSend(m)
}

func getCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}