# rpc 

查看注册服务源码

```go
// 使用 rpc.Register(userServer) 注册服务时，默认使用了 DefaultServer 对象
// 看到这里这种实现方式和http.DefaultServeMux如出一辙

// Register publishes the receiver's methods in the [DefaultServer].
func Register(rcvr any) error { return DefaultServer.Register(rcvr) }

// DefaultServer 通过NewServer() 创建，因此也是可以自己new一个
var DefaultServer = NewServer()

```

当调rpc.Register注册完后，然后就是启动tcp服务，等待tcp的连接，
当有连接进来，把链接交由 rpc 处理 `rpc.ServeConn(conn)`,
启动tcp服务，这里是采用net包，而不是rpc包。因此可以理解 rpc 业务处理能力是和网络区分开的。
rpc 包并没有 net.Listen 这种方法。


```go
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("监听失败: ", err)
	}

	log.Println("服务启动成功 ", addr)

	// 接受连接
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("接受客户端连接失败: ", err)
			continue
		}

        // 就是 DefaultServer
        //源码：
        /**
            func ServeConn(conn io.ReadWriteCloser) {
                DefaultServer.ServeConn(conn)
            }
        */
		rpc.ServeConn(conn)
	}
```

客户端如何调用rpc服务？
从服务端代码看出，将实现业务的对象注册到rpc上，启动tcp服务，接受链接，把链接交由rpc处理。

net tcp 客户端链接是使用 net.Dial() 来拨号连接的。rpc其实也有Dial方法，它所建立tcp链接
就有rpc处理。

```go
// 使用 rpc.Dial 来拨号，建立链接
client, err := rpc.Dial("tcp", "localhost:9007")
	if err != nil {
		log.Fatal("建立连接失败: ", err)
	}
	defer client.Close()

	var (
		req  = GetUserReq{Id: "3"}
		resp GetUserResp
	)

	err = client.Call("UserService.GetUser", req, &resp)
	if err != nil {
		log.Println("请求失败: ", err)
		return
	}

	log.Println(resp)
```

提出问题

rpc 注册的时候做了什么？

client.Call 的时候，如何查找对应的方法？

是所有的对象都可以注册到rpc上吗？
并不是的，必须要满足几个条件，不然是不符合的，在 Register 时就会返回错误了。
- 必须是导出类型（简单来说，就是对象首字母要大写）
- 类型必须有导出的方法
- 导出的方法必须有两个参数，第一个参数必须是指针类型，返回值就绑定到第二个参数上
- 注册的时候，还必须传入指针类型
- 都一一测试过，好的代码注释都写的好清楚。
``` go
// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//   - exported method of exported type
//   - two arguments, both of exported type
//   - the second argument is a pointer
//   - one return value, of type error
//
// It returns an error if the receiver is not an exported type or has
// 客户端调用方式 "Type.Method"
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (server *Server) Register(rcvr any) error {
	return server.register(rcvr, "", false)
}

```

注册源码

``` go
func (server *Server) register(rcvr any, name string, useName bool) error {
    // new service 实例，保存rcvr(注册对象指针)的
    // reflect.Type 和 reflect.Value
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := name
	if !useName {
		sname = reflect.Indirect(s.rcvr).Type().Name()
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(s)
		return errors.New(s)
	}
	if !useName && !token.IsExported(sname) {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	s.name = sname

	// Install the methods（注册方法）
    // 得到的 method 是 map[string]*methodType
    // suitableMethods 查询满足rpc的方法并返回
	s.method = suitableMethods(s.typ, logRegisterError)

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PointerTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Print(str)
		return errors.New(str)
	}

    // LoadOrStore 如果key不存在，就存储
	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}


// suitableMethods returns suitable Rpc methods of typ. It will log
// errors if logErr is true.
func suitableMethods(typ reflect.Type, logErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if !method.IsExported() {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			if logErr {
				log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if logErr {
				log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Pointer {
			if logErr {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if logErr {
				log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if logErr {
				log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if logErr {
				log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}
```


