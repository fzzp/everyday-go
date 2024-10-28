package main

import (
	"log"
	"net/rpc"
)

type (
	GetUserReq struct {
		Id string `json:"id"`
	}

	GetUserResp struct {
		Id    string
		Name  string
		Phone string
	}
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:9007")
	if err != nil {
		log.Fatal("建立连接失败: ", err)
	}
	defer client.Close()

	var (
		req  = GetUserReq{Id: "3"}
		resp GetUserResp
	)

	err = client.Call("UserService.GetUser", req, &resp)
	if err != nil {
		log.Println("请求失败: ", err)
		return
	}

	log.Println(resp)
}
