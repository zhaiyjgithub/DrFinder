package Utils

import (
	"DrFinder/src/response"
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

func ValidateParam(ctx iris.Context, validate *validator.Validate, param interface{}) error {
	if err:= ctx.ReadJSON(param); err != nil {
		response.Fail(ctx, response.Err, response.ParamErr, nil)
		return  err
	}

	err:= validate.Struct(param)
	if err != nil {
		fmt.Println(err.Error())
		response.Fail(ctx, response.Err, err.Error(), nil)

		return err
	}

	return nil
}
