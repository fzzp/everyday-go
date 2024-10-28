package main

type User struct {
	Id    string
	Name  string
	Phone string
}

var users = map[string]*User{
	"1": {
		Id:    "1",
		Name:  "Cat",
		Phone: "13633333333",
	},
	"2": {
		Id:    "2",
		Name:  "Dog",
		Phone: "18777777777",
	},
	"3": {
		Id:    "3",
		Name:  "Fish",
		Phone: "15634243232",
	},
}
