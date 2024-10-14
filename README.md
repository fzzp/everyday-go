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


