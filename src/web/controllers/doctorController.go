package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type DoctorController struct {
	Ctx iris.Context //自动绑定上下文
}

type User struct {
	Name string
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle("POST", "/GetDoctorById", "GetDoctorById")
}

func (c *DoctorController) GetDoctorById() string {
	var user User

	if err:= c.Ctx.ReadJSON(&user); err != nil {
		fmt.Println("failed")
	}else {

		fmt.Println(user.Name)
	}

	return user.Name
}


