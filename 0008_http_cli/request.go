package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RequestConfig 请求
type RequestConfig struct {
	Method   string            // 请求方法    -X
	Data     string            // 发送的数据  -d
	Headers  map[string]string // 请求Header -H
	Output   string            // 请求输出目标 -o
	ShowConf bool              // 打印请求配置 -c
	URL      string            // 请求地址
}

func (conf *RequestConfig) MakeRequest() error {
	// 是否输出配置
	if conf.ShowConf {
		conf.LogConfig()
	}

	// 1. 实例化一个http发送请求的客户端
	client := &http.Client{}

	buf, err := json.Marshal(conf.Data)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(buf)
	fmt.Println(string(buf))

	// 2. 创建一个请求，提供给client发送
	req, err := http.NewRequest(conf.Method, conf.URL, body)
	if err != nil {
		return err
	}

	// 添加headers
	for key, val := range conf.Headers {
		req.Header.Set(key, val)
	}

	// 3. 发送请求
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// data := make([]byte, 1024)
	// n, _ := res.Request.Body.Read(data)
	// fmt.Println("data ->>>", string(data[:n]))

	databuf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("请求成功！")
	fmt.Println(string(databuf))

	// TODO: -o 输出结果到指定文件

	//FIXME: JSON 解析不对？

	return nil
}

func (conf *RequestConfig) LogConfig() {
	buf, _ := json.MarshalIndent(conf, "", " ")
	fmt.Println(string(buf))
}
