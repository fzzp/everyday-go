package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 定义任务接口
type Task interface {
	// 执行任务的方法
	HandleTask()
}

// MessageTask 定义一个消息的任务
type MessageTask struct {
	ID      int    // 消息id
	UserID  int    // 消息接受者id
	Content string // 消息内容
}

// HandleTask 实现 Task 接口
func (mt *MessageTask) HandleTask() {
	time.Sleep(time.Second * 2) // 假设业务繁忙
	fmt.Printf("消息[%d]发发出，接受者[%d] -内容[%s]\n", mt.ID, mt.UserID, mt.Content)
}

// WorkerPool 工作池，控制并发数量和执行任务
type WorkerPool struct {
	Tasks       []Task    // 任务列表
	concurrency int       // 并发数量，控制go协程数量
	tasksChan   chan Task // 任务通道
	wg          sync.WaitGroup
}

// NewWorkerPool 创建一个WorkerPool实例
func NewWorkerPool(tasks []Task, concurrency int) *WorkerPool {
	return &WorkerPool{
		Tasks:       tasks,
		concurrency: concurrency,
		tasksChan:   make(chan Task, len(tasks)),
	}
}

// worker 工作者，真正执行任务的方法
func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.HandleTask()
		wp.wg.Done()
	}
}

// Start 启动任务
// 会初始化任务通道，创建Go协程
func (wp *WorkerPool) Start() {
	// 初始化任务通道
	if wp.tasksChan == nil {
		wp.tasksChan = make(chan Task, len(wp.Tasks))
	}

	// 控制并发数量，先启动工作者，监听通道channel
	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}

	// 发送任务到任务channel
	wp.wg.Add(len(wp.Tasks))

	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}

	// 关闭channel,释放资源
	close(wp.tasksChan)

	// 等待工作者（worker）完成任务
	wp.wg.Wait()
}