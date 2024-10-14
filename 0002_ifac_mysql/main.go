package main

// 这个例子探讨一下mysql+interface实践

// 是否听过这样一句话：
// “Accept interfaces, return structs” in Go
// 在Go中 "接受接口，返回结构体" 是最佳实践

import (
	"database/sql"
	"encoding/json"
	"everyday-go/0002_ifac_mysql/db"
	"everyday-go/0002_ifac_mysql/model"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var dbPwsd string
	flag.StringVar(&dbPwsd, "dbPwsd", "", "数据库密码")
	flag.Parse()

	conn := openDB(dbPwsd)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(10 * time.Minute)

	repo := db.NewRepository(conn)

	getProductById(repo, 1)

	// testInsertWithCategoryTx(repo)

	// insertProduct(repo)
	// insertOneCategory(repo)
	// insertManyCategory(repo)

}

func getProductById(repo *db.Repository, id int) {
	goods, err := repo.Products.QueryRow(id)
	if err != nil {
		log.Println(err)
		return
	}

	buf, _ := json.MarshalIndent(&goods, "", " ")

	fmt.Println(string(buf))
}

func testInsertWithCategoryTx(repo *db.Repository) {
	p := model.Product{ProductName: "IPhone 18 pro", ProductPrice: 666}

	list := make([]model.Category, 0, 3)
	list = append(list, model.Category{CategoryName: "3C"})
	list = append(list, model.Category{CategoryName: "手机"})
	list = append(list, model.Category{CategoryName: "数码"})

	err := repo.Products.InsertWithCategoryTx(p, list)
	if err != nil {
		log.Println(err)
	}
}

func insertProduct(repo *db.Repository) {
	p := model.Product{ProductName: "xiaomi 14 pro", ProductPrice: 999}
	newID, err := repo.Products.Insert(p)
	if err != nil {
		log.Println(err)
	}
	log.Println("newID: ", newID)
}

func insertOneCategory(repo *db.Repository) {
	category := model.Category{CategoryName: "笔记本"}
	newID, err := repo.Category.Insert(category)
	if err != nil {
		log.Println(err)
	}
	log.Println("newID: ", newID)
}

func insertManyCategory(repo *db.Repository) {
	list := make([]model.Category, 0, 3)
	list = append(list, model.Category{CategoryName: "3C"})
	list = append(list, model.Category{CategoryName: "手机"})
	list = append(list, model.Category{CategoryName: "数码"})

	err := repo.Category.InsertMany(list)
	if err != nil {
		log.Println(err)
	}
}

func openDB(dbPwsd string) *sql.DB {
	var (
		username = "root"
		dbname   = "db_everyday_go"
	)

	// 链接mysql并设置时区跟随本地系统
	dsn := fmt.Sprintf(
		"%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&loc=Local&parseTime=True",
		username,
		dbPwsd,
		dbname,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
