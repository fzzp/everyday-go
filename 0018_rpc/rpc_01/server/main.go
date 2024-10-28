package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type UserService struct{}

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

// GetUser req 是请求入参，resp 是请求结果，通过指针返回
func (*UserService) GetUser(req GetUserReq, resp *GetUserResp) error {
	if u, ok := users[req.Id]; ok {
		*resp = GetUserResp{
			Id:    u.Id,
			Name:  u.Name,
			Phone: u.Phone,
		}

		return nil
	}

	return errors.New("没有找到用户")
}

func main() {
	// 创建服务
	userServer := new(UserService)

	// 服务注册到rpc
	if err := rpc.Register(userServer); err != nil {
		log.Fatal("注册服务失败: ", err)
	}

	addr := ":9007"

	// 监听
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("监听失败: ", err)
	}

	log.Println("服务启动成功 ", addr)

	// 接受连接
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("接受客户端连接失败: ", err)
			continue
		}

		rpc.ServeConn(conn)
	}
}
