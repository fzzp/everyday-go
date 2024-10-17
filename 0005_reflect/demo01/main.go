package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 基础变量
	var s = "hi"

	fmt.Println("s type: ", reflect.TypeOf(s))   // string
	fmt.Println("s value: ", reflect.ValueOf(s)) // hi

	t := reflect.TypeOf(s)
	fmt.Println(t.Name()) // string

	v := reflect.ValueOf(s)
	fmt.Println(v.Type())
	fmt.Println(v.String())

	var a interface{}
	a = 10
	a = "555"
	a = true
	fmt.Println(a)
}
