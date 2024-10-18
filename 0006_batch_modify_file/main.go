package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// findVtt 查找vtt字幕文件
func findVtt() []string {
	vtts, err := filepath.Glob("**/*.vtt")
	if err != nil {
		panic(err)
	}
	return vtts
}

type VttTask string

func (vtt VttTask) HandleTask() {
	file, err := os.OpenFile(string(vtt), os.O_RDWR, 0644)
	if err != nil {
		log.Printf(string(vtt)+" 打开文件错误 %v : %s\n", vtt, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var delLine1 = "Kind: captions"
	var delLine2 = "Language: en"

	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		flag := true
		if line == delLine1 {
			line = strings.ReplaceAll(line, delLine1, "")
			flag = false
		}
		if line == delLine2 {
			line = strings.ReplaceAll(line, delLine2, "")
			flag = false
		}

		if flag {
			lines = append(lines, line)
		}
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println(string(vtt)+" 扫描文件错误: ", err)
		return
	}

	var content bytes.Buffer
	for _, line := range lines {
		content.WriteString(line + "\n")
	}

	err = os.WriteFile(string(vtt), content.Bytes(), 0644)
	if err != nil {
		fmt.Println(string(vtt)+" 写入文件错误: ", err)
		return
	}

	// fmt.Println(string(vtt) + " 修改成功！")
}

func main() {
	vtts := findVtt()
	if len(vtts) == 0 {
		log.Println("没有文件可处理")
		return
	}

	var tasks []Task
	for _, vtt := range vtts {
		tasks = append(tasks, VttTask(vtt))
	}

	wp := NewWorkerPool(tasks, 20)
	wp.Start()
}
