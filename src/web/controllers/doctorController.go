package controllers

import (
	"DrFinder/src/models/doctor"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"time"
)

type DoctorController struct {
	Ctx iris.Context
	Service service.DoctorService
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("POST", "/AddDoctor", "AddDoctor")
}

func (c *DoctorController) AddDoctor() error {
	type doctorParam struct {
		Name string
	}

	var param doctorParam

	if err:= c.Ctx.ReadJSON(&param); err == nil {

	}else {

	}

	newDoctor:= &doctor.Doctor{
		Npi: 1316960271,
		LastName: "MASSEY",
		FirstName: "WILLIAM",
		MiddleName: "A",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: "",
		NamePrefix: "DR.",
		Gender: "M",
		MailingAddress: "504 BROOKWOOD BLVD SUITE 100",
		MailingCity: "BIRMINGHAM",
		MailingState: "AL",
		MailingZip: "352096802",
		MailingPhone: "2058719661",
		MailingFax: "2058701621",
		BusinessAddress: "100 PILOT MEDICAL DR SUITE 100",
		BusinessCity: "BIRMINGHAM",
		BusinessState: "AL",
		BusinessZip: "352353411",
		BusinessPhone: "2058548084",
		BusinessFax: "2058159341",
		Specialty: "Allergy & Immunology",
	}

	err:= c.Service.Add(newDoctor)

	return err
}

