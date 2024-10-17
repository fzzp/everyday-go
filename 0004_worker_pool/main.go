package main

import (
	"math/rand"
	"strings"
	"time"
)

var (
	myRand     *rand.Rand
	characters = "QWERTYUIOPASDLZMqwiopasdfghjklzxcvbnm"
)

func main() {
	var tasksNum = 20   // 任务数
	var concurrency = 5 // go 协程数
	tasks := make([]Task, tasksNum)
	for i := 0; i < tasksNum; i++ {
		userId := myRand.Intn(10000)
		content := RandomString(10)
		tasks[i] = &MessageTask{ID: i + 1, UserID: userId, Content: content}
	}

	wp := NewWorkerPool(tasks, concurrency)
	wp.Start()
}

func init() {
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(n int) string {
	var s strings.Builder

	for i := 0; i < n; i++ {
		b := characters[myRand.Intn(len(characters))]
		s.WriteByte(b)
	}

	return s.String()
}
