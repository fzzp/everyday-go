package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

/*
解析命令参数：

➜  0008_http_cli git:(master) ✗ go run *.go  -X POST -o -c index.json -H '{"Content-Type":"application/json"}' -d '{"pageInt": 10}'  http://baidu.com\?a\=100
{
 "Method": "POST",
 "Data": "{\"pageInt\": 10}",
 "Headers": {
  "Content-Type": "application/json"
 },
 "Output": "index.json",
 "ShowConf": true,
 "URL": "http://baidu.com?a=100"
}
➜  0008_http_cli git:(master) ✗

*/

// ./bin/app -X POST -H '{"Content-Type":"application/json"}' -d '{"email":"lightsaid@163.com", "password":"abc123"}' -c  http://localhost:8901/v1/auth/login

func main() {
	// NOTE: 这里其实也可以自己解析参数也行，flag.Parse() 也是用os.Args来解析的
	// os.Args 返回的就是所有命令行参数
	// log.Println(os.Args[1:])

	// -c 参数在这里后面不需要跟随参数，默认用户输出-c就是true了，
	// 因此，-c 交由 flag 处理是不正确的
	// 而flag.Parse()用的os.Args全局变量
	// 所有，在flag.Parse()执行之前需要剔除-c

	args, c := findKeyByArgs(os.Args, "-c")
	os.Args = args

	r := RequestConfig{}
	r.URL = os.Args[len(os.Args)-1]

	var headerJson string
	flag.StringVar(&r.Method, "X", "GET", "请求方法")
	flag.StringVar(&r.Data, "d", "", "请求body参数")
	flag.StringVar(&headerJson, "H", "", "请求头JSON String")
	flag.StringVar(&r.Output, "o", "", "请求结果输出")
	flag.StringVar(&r.Output, "c", "", "是否输出请求配置") // 此处仅为占位，不然输入 -c 参数会报错

	flag.Parse()
	r.ShowConf = c

	if len(headerJson) > 0 {
		if err := json.Unmarshal([]byte(headerJson), &r.Headers); err != nil {
			fmt.Println("Header格式不正确")
			return
		}
	}

	// buf, _ := json.MarshalIndent(r, "", " ")
	// fmt.Println(string(buf))

	if ok := ValidateURL(r.URL); !ok {
		fmt.Println("请求地址不正确")
		return
	}

	if err := r.MakeRequest(); err != nil {
		log.Fatal(err)
	}
}
