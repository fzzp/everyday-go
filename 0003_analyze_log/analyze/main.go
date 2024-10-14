package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	goods := make(map[int]int, 2)
	goods[11] = 0
	goods[22] = 0

	start := time.Now()
	analyzeLog("access.log", &goods)
	fmt.Printf("%s\n", time.Since(start))

	fmt.Println(goods)

	gzCompressLog("access.log", "access.log.gz")
}

func analyzeLog(logfile string, m *map[int]int) {
	mutex := sync.RWMutex{}
	f, err := os.Open(logfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "GET - /product/11") {
			mutex.Lock()
			(*m)[11] = (*m)[11] + 1
			mutex.Unlock()

		}
		if strings.Contains(line, "GET - /product/22") {
			mutex.Lock()
			(*m)[22] = (*m)[22] + 1
			mutex.Unlock()
		}
	}
}

// gzCompressLog 压缩日志 src输入，dst输出
func gzCompressLog(src, dst string) error {
	// 打开原文件
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建压缩输出文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 创建压缩写入端
	writer := gzip.NewWriter(out)
	defer writer.Close()

	writer.Name = src

	info, err := file.Stat()
	if err == nil {
		// 复制一下修改时间
		writer.ModTime = info.ModTime()
	}

	// 使用 io.Copy 复制
	if _, err := io.Copy(writer, file); err != nil {
		os.Remove(dst) // 如果复制错误，删除一下
		return err
	}

	return nil
}
