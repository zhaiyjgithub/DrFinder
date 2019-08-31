package main

import (
	"DrFinder/src/conf"
	"fmt"
	"time"
)

const format = "2006-01-02 15:04:05"

func main()  {
	t, _ := time.Parse(conf.TimeFormat, "2019-08-31 12:37:25")
	t1, _ := time.Parse(conf.TimeFormat, "2019-08-31 12:55:15")

	stamp := t.Unix()
	stamp1 := t1.Unix()

	diff := stamp1 - stamp

	fmt.Println(diff)
}
