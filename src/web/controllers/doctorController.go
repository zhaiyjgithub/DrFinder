package controllers

import (
	"DrFinder/src/models/doctor"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type DoctorController struct {
	Ctx iris.Context
	Service service.DoctorService
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("POST", "/AddDoctor", "AddDoctor")
}

func (c *DoctorController) addDoctor()  {

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

	}

	err:= c.Service.Add(newDoctor)

	return err
}

