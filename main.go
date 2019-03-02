// hello_world project main.go
package main

import (
	//	"app"
	"fmt"
)

var app Application

func main() {
	fmt.Println("Hello World!")
	if app.init() {
		app.run()
	}

	//gpio.Open(1)
}
