package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

type RegisterController struct {
	Ctx iris.Context
	Service service.UserService
}

type Code struct {
	Email string `json:"email"`
	TimeStamp string `json:"date"`
	Value string `json:"value"`
}

var signValidate *validator.Validate

func (c *RegisterController) BeforeActivation(b mvc.BeforeActivation)  {
	signValidate = validator.New()
	b.Handle(iris.MethodPost, Utils.SendVerificationCode, "SendVerificationCode")
}

func (c *RegisterController) SendVerificationCode() {
	type Param struct {
		Email string `validate:"email"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, signValidate, &param)

	if err != nil  {
		return
	}

	user := c.Service.GetUserByEmail(param.Email)

	if user != nil {
		response.Fail(c.Ctx, response.Err, "This email has been registered", nil)
	}else {
		v := getCode()

		var code = &Code{
			Email: param.Email,
			TimeStamp:  time.Now().Format("20060102150405"),
			Value: v,
		}

		cb, _ := json.Marshal(code)
		err = dataSource.Save(param.Email, string(cb))

		err := sendEmail(param.Email, code.Value)
		if err != nil {
			response.Fail(c.Ctx, response.Err, "send email fail", nil)
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