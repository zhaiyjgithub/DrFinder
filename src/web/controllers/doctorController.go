package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
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
		Name string `validate:"gt=0"`
		Address string `validate:"gt=0"`
	}

	var param doctorParam

	if err:= c.Ctx.ReadJSON(&param); err != nil {
		response.Fail(c.Ctx, response.Err, response.ParamErr, nil)
		return
	}

	err:= validate.Struct(&param)
	if err != nil {
		fmt.Println(err.Error())
	}

	//newDoctor:= &doctor.Doctor{
	//	Npi: 1316960271,
	//	LastName: "MASSEY",
	//	FirstName: "WILLIAM",
	//	MiddleName: "A",
	//	CreatedAt: time.Now(),
	//	UpdatedAt: time.Now(),
	//	Name: "",
	//	NamePrefix: "DR.",
	//	Gender: "M",
	//	MailingAddress: "504 BROOKWOOD BLVD SUITE 100",
	//	MailingCity: "BIRMINGHAM",
	//	MailingState: "AL",
	//	MailingZip: "352096802",
	//	MailingPhone: "2058719661",
	//	MailingFax: "2058701621",
	//	BusinessAddress: "100 PILOT MEDICAL DR SUITE 100",
	//	BusinessCity: "BIRMINGHAM",
	//	BusinessState: "AL",
	//	BusinessZip: "352353411",
	//	BusinessPhone: "2058548084",
	//	BusinessFax: "2058159341",
	//	Specialty: "Allergy & Immunology",
	//}
	//
	//ok := c.Service.Add(newDoctor)
	//
	//if ok == true {
	//	response.Success(c.Ctx, response.Successful, nil)
	//}else {
	//	response.Fail(c.Ctx, response.Err, "", nil)
	//}
}
