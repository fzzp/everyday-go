package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Author struct {
	Name string
}

type Book struct {
	ID    int
	Title string
	Author
}

func (b Book) Show() {
	buf, _ := json.MarshalIndent(b, "", " ")
	fmt.Println(string(buf))
}

func (b *Book) Update(title string) {
	b.Title = title
	b.Show()
}

// test1 通过 Interface 获取对应的数值
func test1() {
	b := Book{ID: 1, Title: "Golang"}

	v := reflect.ValueOf(b)

	fmt.Println(v.Interface()) // {1 Golang}

	// Interface 断言
	b2, ok := v.Interface().(Book)
	if !ok {
		fmt.Println(ok)
		return
	}

	fmt.Println(b2.ID, " - ", b2.Title) // 1  -  Golang
}

// test2 获取字段和对应的值
func test2() {
	author := Author{Name: "zhangsan"}
	book := Book{ID: 1, Title: "Golang", Author: author}
	rType := reflect.TypeOf(book)
	fmt.Println(rType.Name()) // Book
	fmt.Println(rType.Kind()) // struct （kind种类，基本数据类型）

	rValue := reflect.ValueOf(book)

	fmt.Println(rType.NumField())

	// 获取 field 字段
	for i := 0; i < rType.NumField(); i++ {
		// 获取到字段信息
		field := rType.Field(i)

		// 字段名称 - 类型 - 种类
		// ID - int - int
		// Title - string - string
		// Author - main.Author - struct
		fmt.Println(field.Name, "-", field.Type, "-", field.Type.Kind())

		// 获取每个字段对应的值
		// [ID]=1
		// [Title]=Golang
		// [Author]={zhangsan}  // NOTE: 如果还需要获取Author对象，继续对它进行反射
		val := rValue.Field(i).Interface()
		fmt.Printf("[%s]=%v\n", field.Name, val)
	}
}

// test03 通过反射设置值
// 设置值，必须传递指针给reflect.ValueOf(ptr)，否则会panic
func test03() {
	author := Author{Name: "zhangsan"}
	book := Book{ID: 1, Title: "Golang", Author: author}

	// NOTE:
	// 对于指针的反射，基本都需要使用 Elem() 方法获取指针指向的值
	// 如果传递的值不是指针，使用Elem()方法是会panic
	// 因此在合适的位置需要使用recover()捕获panic
	rValue := reflect.ValueOf(&book).Elem()

	// 查找字段, 可以通过名字获取或者下标获取
	title := rValue.FieldByName("Title")
	// 上面字段如果没有，CanSet 肯定就是false
	if title.CanSet() {
		title.SetString("Docker")
	}

	// 最终修改了title
	fmt.Println(book) // {1 Docker {zhangsan}}
}

// test04 通过反射调用方法
func test04() {
	author := Author{Name: "zhangsan"}
	book := Book{ID: 1, Title: "Golang", Author: author}

	// 两种方式都可以
	// rVal := reflect.ValueOf(&book).Elem()
	rVal := reflect.ValueOf(book)

	show := rVal.MethodByName("Show")
	// 如果方法不存在，就不用往下执行了
	if !show.IsValid() || show.Kind() != reflect.Func {
		fmt.Println("show方法不存在")
	}
	// kind: func, type:func()
	fmt.Printf("kind: %s, type:%s\n", show.Kind(), show.Type())

	// 执行Show方法，没有参数就传nil
	show.Call(nil)

	// NOTE: 当方法绑定在指针上的，这种方式获取不到的
	// update := rVal.MethodByName("Update")

	update := rVal.MethodByName("Update")
	if !update.IsValid() { // false 是无效， 有效返回true
		fmt.Println("获取不到指针的方法")
	} else {
		fmt.Printf("kind: %s, type:%s\n", update.Kind(), update.Type())
	}

	// NOTE: 获取指针值的方法, 参考 https://go.dev/play/p/iiUrB961tUa
	// 获取指针的方法，使用指针作为reflect.ValueOf(&book)入参
	// 但是不能转换成真实的对象了，因为方法就绑定在指针上，因此
	// reflect.ValueOf 必须传指针
	bookVal := reflect.ValueOf(&book)
	updateFunc := bookVal.MethodByName("Update")
	if updateFunc.IsValid() {
		arg := []reflect.Value{reflect.ValueOf("新标题")}
		updateFunc.Call(arg)
	} else {
		fmt.Println("updateFunc 不存在")
	}

	// Show 也是可以获取的
	showFunc := bookVal.MethodByName("Show")
	showFunc.Call(nil)

	// NOTE: 总结：反射调用方法，全部使用reflect.ValueOf(&book)传
	// 指针即可，不要使用Elem方法
}

// test05 通过反射创建结构体对象
func test05() {
	author := Author{Name: "zhangsan"}
	book := Book{ID: 1, Title: "Golang", Author: author}

	bookType := reflect.TypeOf(book)

	// NOTE: reflect.New 返回的是一个对象的指针 reflect.Value
	// 这里需要使用 Elem 获取值实际的value，不然下面代码会报错
	bookValue := reflect.New(bookType).Elem()

	// 这里不做过多判断直接设置
	bookValue.FieldByName("ID").SetInt(22)
	bookValue.FieldByName("Title").SetString("TypeScript")

	book2 := bookValue.Interface().(Book)

	// 两个对象就没有关联了
	fmt.Println(book.ID, book.Title)
	fmt.Println(book2.ID, book2.Title)

	// NOTE: 如果要创建map或者slice 就是需要用到下面的方法了
	// reflect.MakeSlice
	// reflect.MakeMapWithSize(

	// 诸如此类此类，多看文档，先复习到这里
}

func main() {
	// test1()
	// test2()
	// test03()
	// test04()
	test05()
}
