package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
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
	b.Handle(iris.MethodPost, Utils.GetDoctorById, "GetDoctorById")
	b.Handle(iris.MethodPost, Utils.SearchDoctor, "SearchDoctor")
	b.Handle(iris.MethodPost, Utils.UpdateDoctorById, "UpdateDoctorById")
	b.Handle(iris.MethodPost, Utils.DeleteDoctorById, "DeleteDoctorById")
}

func (c *DoctorController) AddDoctor() {
	type Param struct {
		Npi             int64  `validate:"gt=0,numeric"`
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

	var param Param

	err:= Utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	newDoctor:= &models.Doctor{
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

func (c *DoctorController)GetDoctorById() {
	type Param struct {
		DoctorId int `validate:"gt=0,numeric"`
	}

	var param Param

	err:= Utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	doctor := c.Service.GetDoctorById(param.DoctorId)

	if  doctor != nil {
		response.Success(c.Ctx, response.Successful, doctor)
	}else {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}
}

func (c *DoctorController) GetDoctorBySpecialty(specialty string)  {
	type Param struct {
		Specialty string `validate:"gt=0"`
	}

	var param Param
	err:= Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	doctor:= c.Service.GetDoctorBySpecialty(param.Specialty)

	if doctor != nil {
		response.Success(c.Ctx, response.Successful, doctor)
	}else {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}
}

func (c *DoctorController) SearchDoctor()  {
	type Param struct {
		FirstName string
		LastName string
		Specialty string
		Gender string `validate:"len=1"`
		City string
	}

	var param Param
	err:= Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var doctor models.Doctor
	doctor.FirstName = param.FirstName
	doctor.LastName = param.LastName
	doctor.Specialty = param.Specialty
	doctor.Gender = param.Gender
	doctor.BusinessCity = param.City

	doctors:= c.Service.SearchDoctor(&doctor)

	response.Success(c.Ctx, response.Successful, doctors)
}

func (c *DoctorController) UpdateDoctorById() {
	type Param struct {
		ID        int
		FirstName string
	}

	var param Param
	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var doctor models.Doctor
	doctor.FirstName = param.FirstName
	doctor.ID = param.ID

	err = c.Service.UpdateDoctorById(&doctor)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "update failed", nil)
	} else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *DoctorController) DeleteDoctorById()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err:= Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	ok:= c.Service.DeleteDoctorById(param.ID)

	if ok {
		response.Success(c.Ctx, response.Successful, nil )
	}else {
		response.Fail(c.Ctx, response.Err, "delete failed", nil)
	}
}
