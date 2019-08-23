package controllers

import "github.com/kataras/iris/mvc"

type DoctorController struct {

}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("GET", "/GetDoctorById", "GetDoctorById")
}

func (c *DoctorController) GetDoctorById() string {
	return "doctor not found"
}

