# Everyday Go 

每天写点Go代码。

## 探索Go语言

根据实际需求去写Go代码。

### 0001_io_demo 

如何将基础数据（string/struct ...）转为实现 io.Reader 或 io.Writer 接口的数据

```go
bytes.NewReader
strings.NewReader
bytes.Buffer
```

### 0002_ifac_mysql 

运行 `go run main.go -dbPwsd=本地数据库密码`

这个例子使用原生SQL优雅实现数据库CRUD操作。

这里以一个简单的例子，包括单行/多行插入，join 表查询，事务查询。

我的选择：
sqlc > sqlx > gorm 


### 0003_analyze_log

这个例子介绍日志分析和压缩。

首先在genlog创建一个简单的web服务，并访问生成访问日志。
到genlog执行`go run *.go`生成10w+行日志，将logs移动到analyze目录。

去到 analyze 执行分析，读取日志，分析商品id为11和22的分别访问了多少次。

最终结果如下，不得不说使用 `bufio.NewScanner()`扫描文件真的快，22M大小日志，25w+行，扫描完成才用了43ms；压缩成629K。

```bash
➜  analyze git:(master) ✗ go run main.go
43.273007ms
map[11:63729 22:56480]
➜  analyze git:(master) ✗ ll
total 45704
-rw-r--r--@ 1 lightsaid  staff    22M Oct 14 22:35 access.log
-rw-r--r--@ 1 lightsaid  staff   629K Oct 14 23:07 access.log.gz
-rw-r--r--@ 1 lightsaid  staff   1.5K Oct 14 23:07 main.go
➜  analyze git:(master) ✗ 
```


### 0004_worker_pool

Worker Pool（工作池/协程池）模式是一种常见的并发设计模式，主要用于控制并发任务的数量，提高系统性能，以及更有效地管理系统资源。

实现一个简单的Worker-Pool模式；

工作上肯定是用 [ants](https://github.com/panjf2000/ants) 控制并发数量。
[ants文档](https://github.com/panjf2000/ants/blob/v2.10.0/README_ZH.md)
[ants解读](https://mp.weixin.qq.com/s/Uctu_uKHk5oY0EtSZGUvsA)

`简单回顾一下channel：`
    
非缓冲通道：make(chan T)
>>>
    一次发送，一次接受，都是阻塞的

缓冲通道：make(chan T, capacity)
>>>
    发送：缓冲区数据满了，才会阻塞
    接收：缓冲区数据空了，才会阻塞

双向通道：
>>>
    chan <- data 发送数据，写入
    data <- chan 接收数据，读取

定向（单向）通道：只能接受或者发送数据
>>> 
    chan <- T,  只写
    <- chan T,  只读

定向通道定义：
```go
ch1 := make(chan<- int) // 只写
ch2 := make(<-chan int) // 只读
```

### 0005_reflect 

复盘反射。

reflect.TypeOf -> reflect.Type

reflect.Type 
> 是和类型相关的操作，如：查看字段信息(个数、名字、tag等)，根据类型创建新的对象


reflect.ValueOf -> reflect.Value

reflect.Value
> 是和值相关的操作，获取值，设置新的值，调用方法

### 0006_batch_modify_file

今天下载了一份老外的教程，配有英文字幕，但是播放器识别不了该字幕，原因是字幕开头有两行多余影响了。
将这两行删除即可，但是将近200多个文件，手动一个个去打开删除，不但效率低，还恐有纰漏。

于是就写这个程序来删除吧。

```txt
WEBVTT

Kind: captions  (删除此行)

Language: en  (删除此行)

00:00:01.972 --> 00:00:07.266

```

go build -o=app 后，添加可执行权 chmod +x app，再移到课程里，执行即可。
