/*
Copyright © 2023 Killua<captainchengjie@gmail.com>

*/
package main

import (
	"fmt"
	// "ginBlog/api"
	"ginBlog/cmd"
)

func main() {
	fmt.Println("start")
	cmd.Execute()
	// data := <-api.A
	// fmt.Println(data)
	fmt.Println("end")
}
