package main

import (
	"fmt"
	"time"
)

func main() {
	s := time.Now()
	time.Sleep(time.Second * 5)
	e := time.Since(s)
	fmt.Println(e.Milliseconds())
}
