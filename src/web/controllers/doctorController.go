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

	_ "github.com/kataras/iris/sessions/sessiondb/boltdb"
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
	b.Handle(iris.MethodPost, Utils.SearchDoctorByPage, "SearchDoctorByPage")
}

func (c *DoctorController) AddDoctor() {
	type Param struct {
		Npi          int  `validate:"gt=0,numeric"`
		LastName     string `validate:"gt=0"`
		FirstName    string `validate:"gt=0"`
		MiddleName   string
		FullName     string `validate:"gt=0"`
		NamePrefix   string
		Credential   string `validate:"gt=0"`
		Gender       string `validate:"len=1"`
		Address      string `validate:"gt=0"`
		City         string `validate:"gt=0"`
		State        string `validate:"gt=0"`
		Zip          string `validate:"gt=0"`
		Phone        string
		Fax 		 string
		Specialty    string `validate:"gt=0"`
		SubSpecialty string
		JobTitle     string
		Summary      string
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
		FullName: param.FullName,
		NamePrefix: param.NamePrefix,
		Gender: param.Gender,
		Address: param.Address,
		City: param.City,
		State: param.State,
		Zip: param.Zip,
		Phone: param.Phone,
		Fax: param.Fax,
		Specialty: param.Specialty,
		SubSpecialty: param.SubSpecialty,
		JobTitle: param.JobTitle,
		Summary: param.Summary,
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
		Gender string
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
	doctor.City = param.City

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

func (c *DoctorController) SearchDoctorByPage()  {
	type Param struct {
		FirstName string
		LastName string
		Specialty string
		Gender string
		City string
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var doctor models.Doctor
	doctor.FirstName = param.FirstName
	doctor.LastName = param.LastName
	doctor.Specialty = param.Specialty
	doctor.Gender = param.Gender
	doctor.City = param.City

	doctors := c.Service.SearchDoctorByPage(&doctor, param.Page, param.PageSize)

	response.Success(c.Ctx, response.Successful, doctors)
}
