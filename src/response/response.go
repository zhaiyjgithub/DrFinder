package response

import "github.com/kataras/iris"



const Ok = 0
const Err = 1



const Successful = "success"
const ParamErr = "param error"
const NotFound = "Nor found"



func Success(ctx iris.Context, msg string, data interface{})  {
	ctx.JSON(iris.Map{
		"data": data,
		"code": Ok,
		"msg": msg,
	})
}

func Fail(ctx iris.Context, code int, msg string, data interface{})  {
	ctx.JSON(iris.Map{
		"code": code,
		"msg": msg,
		"data": data,
	})
}
