package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// NOTE: 如何将数据转为实现 io.Reader 或 io.Writer 的数据呢？
// 比如现在有一个结构体或者map类型再或者普通的数据类型（int,float,bool,string）
// 这些基础数据都是没有实现io接口的
// 那么如何转换呢？
// 下面例子就是做这个件事情

// bytes.NewReader
// strings.NewReader
// bytes.Buffer

type Book struct {
	ID    int
	Title string
}

func main() {
	var bookJSON = "{\"ID\":7,\"Title\":\"Go\"}"
	reader := strings.NewReader(bookJSON)
	reader.WriteTo(os.Stdout) // 输出到终端

	var book Book
	err := json.Unmarshal([]byte(bookJSON), &book)
	if err != nil {
		panic(err)
	}

	fmt.Println()

	writer := bytes.NewBuffer(make([]byte, len(bookJSON)))
	err = json.NewEncoder(writer).Encode(&book)
	if err != nil {
		panic(err)
	}
	writer.WriteTo(os.Stdout)

	// writer.Bytes()
}
