package main

import (
	"fmt"

	"github.com/chinathaip/chatify/router"
)

func main() {
	fmt.Println("Hello Chatify")

	r := router.New()
	r.Serve(":1337")
}
