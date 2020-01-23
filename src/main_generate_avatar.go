package main

import (
	"fmt"
	"github.com/o1egl/govatar"
)

func main()  {
	fmt.Println("generate avatar")

	err := govatar.GenerateFile(govatar.MALE, "./src/upload/avatar.jpg")

	fmt.Println(err)
}