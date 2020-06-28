package main

import (
	"fmt"
	"sync"
)

func main()  {
	var once sync.Once

	done := make(chan bool)
	for i := 0; i < 10; i++  {
		go func() {
			once.Do(callOnce)
			done <- true
		}()
	}

	for i := 0; i < 10; i ++  {
		<-done
	}

	fmt.Println("finish....")
}

func callOnce()  {
	fmt.Println("only once")
}


