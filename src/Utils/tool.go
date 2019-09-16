package Utils

import (
	"DrFinder/src/conf"
	"DrFinder/src/response"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Level int

const (
		info Level = iota
		warn
	)

func ValidateParam(ctx iris.Context, validate *validator.Validate, param interface{}) error {
	err := ctx.ReadJSON(param)
	if  err != nil {
		response.Fail(ctx, response.Err, response.ParamErr, nil)
		go LogRequest(ctx, param, err.Error(), warn)
		return  err
	}

	err = validate.Struct(param)
	if err != nil {
		fmt.Println(err.Error())
		response.Fail(ctx, response.Err, err.Error(), nil)
		go LogRequest(ctx, param, err.Error(), warn)
		return err
	}

	go LogRequest(ctx, param, "Success", info)

	return nil
}

func LogRequest(ctx iris.Context, v interface{}, error string, level Level)  {
	c := ctx.GetStatusCode()
	r := ctx.GetCurrentRoute()
	b, _ := json.Marshal(v)
	p := string(b)

	log := logrus.New()

	if level == warn {
		log.WithFields(logrus.Fields{
			"route": r,
			"param": p,
			"error": error,
			"statusCode": c,
			"time": time.Now().Format(conf.TimeFormat),
		}).Warnln("API Log")
	}else {
		log.WithFields(logrus.Fields{
			"route": r,
			"param": p,
			"error": error,
			"statusCode": c,
			"time": time.Now().Format(conf.TimeFormat),
		}).Infoln("API Log")
	}
}
