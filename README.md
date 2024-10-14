# Everyday Go 

每天写点Go代码。

## 探索Go语言

根据实际需求去写Go代码。

### 001_io_demo 

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
