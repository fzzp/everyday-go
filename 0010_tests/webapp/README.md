# webapp

go: cannot run *_test.go files (cmd/web/handlers_test.go) 解决方案：
```bash
➜  webapp git:(master) ✗ go run cmd/web/*.go                
go: cannot run *_test.go files (cmd/web/handlers_test.go)
➜  webapp git:(master) ✗ go run `ls cmd/web/*.go | grep -v _test.go`
2024/10/24 23:21:12 Starting server on port  :8631
```