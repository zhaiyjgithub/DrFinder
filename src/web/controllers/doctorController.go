package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models/doctorModel"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

var validate *validator.Validate

type DoctorController struct {
	Ctx     iris.Context
	Service service.DoctorService
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	validate = validator.New()
	b.Handle(iris.MethodPost, Utils.AddDoctor, "AddDoctor")
}

func (c *DoctorController) AddDoctor() {
	type doctorParam struct {
		Npi             int64  `validate:"gt=0"`
		LastName        string `validate:"gt=0"`
		FirstName       string `validate:"gt=0"`
		MiddleName      string
		Name            string
		NamePrefix      string
		Credential      string `validate:"gt=0"`
		Gender          string `validate:"len=1"`
		MailingAddress  string `validate:"gt=0"`
		MailingCity     string `validate:"gt=0"`
		MailingState    string `validate:"gt=0"`
		MailingZip      string `validate:"gt=0"`
		MailingPhone    string
		MailingFax      string
		BusinessAddress string `validate:"gt=0"`
		BusinessCity    string `validate:"gt=0"`
		BusinessState   string `validate:"gt=0"`
		BusinessZip     string `validate:"gt=0"`
		BusinessPhone   string `validate:"gt=0"`
		BusinessFax     string
		Specialty       string `validate:"gt=0"`
	}

	var param doctorParam

	if err:= c.Ctx.ReadJSON(&param); err != nil {
		response.Fail(c.Ctx, response.Err, response.ParamErr, nil)
		return
	}

	err:= validate.Struct(&param)
	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
		return
	}

	newDoctor:= &doctorModel.Doctor{
		Npi: param.Npi,
		LastName: param.LastName,
		FirstName: param.FirstName,
		MiddleName: param.MiddleName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: "",
		NamePrefix: param.NamePrefix,
		Gender: param.Gender,
		MailingAddress: param.MailingAddress,
		MailingCity: param.MailingCity,
		MailingState: param.MailingState,
		MailingZip: param.BusinessZip,
		MailingPhone: param.MailingPhone,
		MailingFax: param.MailingFax,
		BusinessAddress: param.BusinessAddress,
		BusinessCity: param.BusinessCity,
		BusinessState: param.BusinessState,
		BusinessZip: param.BusinessZip,
		BusinessPhone: param.BusinessPhone,
		BusinessFax: param.BusinessFax,
		Specialty: param.Specialty,
	}

	ok := c.Service.Add(newDoctor)

	if ok == true {
		response.Success(c.Ctx, response.Successful, nil)
	}else {
		response.Fail(c.Ctx, response.Err, "", nil)
	}
}

func (c *DoctorController)GetDoctorById(id int) {
	type Param struct {
		ID int
	}

	var param Param

	if err:= c.Ctx.ReadJSON(&param); err != nil {
		response.Fail(c.Ctx, response.Err, response.ParamErr, nil)
		return
	}

	err:= validate.Struct(&param)
	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
		return
	}

	doctor := c.Service.GetDoctorById(id)
	response.Success(c.Ctx, response.Successful,  doctor)
}
