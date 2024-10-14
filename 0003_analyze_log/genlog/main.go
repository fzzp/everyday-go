package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	port   = 4567
	myRand *rand.Rand
	wg     sync.WaitGroup
)

// NOTE: 生成日志
func main() {
	// 打开log file
	logFile := openLogfile()
	defer logFile.Close()

	// 设置slog
	setupSlog(logFile)

	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /product/{pid}", getProductHandler)

	srv := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("r: ", r)
			}
		}()
		// 等待服务启动
		time.Sleep(3 * time.Second)

		// 访问10万次，生成日志再分析
		for i := 0; i < 100_000; i++ {
			time.Sleep(time.Millisecond * 3)
			wg.Add(1)
			go client(i)
		}

		wg.Wait()
		fmt.Println("日志生成完毕～")
	}()

	fmt.Println("start server on ", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// openLogfile 打开日志文件，如果目录或文件不存在就创建它
func openLogfile() *os.File {
	fi, err := os.Stat("logs")
	if err != nil || !fi.IsDir() {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile("logs/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return file
}

// setupSlog slog 日志类型和输出目标
func setupSlog(w io.Writer) {
	logger := slog.New(
		slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	slog.SetDefault(logger)
}

// 客户端，访问商品生成日志
func client(row int) {
	defer wg.Done()
	fmt.Println(row + 1)
	var ids = []int{11, 22}
	id := ids[myRand.Intn(2)]
	url := fmt.Sprintf("http://localhost:%d/product/%d", port, id)
	res, err := http.Get(url)
	if err != nil {
		slog.Error("访问商品错误", "productId", id)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		var p Product
		json.NewDecoder(res.Body).Decode(&p)
		slog.Info("访问成功", "product", p)
	}
}
