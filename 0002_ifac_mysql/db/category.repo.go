package db

import (
	"context"
	"everyday-go/0002_ifac_mysql/model"
	"log"
	"strings"
	"time"
)

type CategoryRepo interface {
	Insert(model.Category) (int64, error)
	InsertMany(list []model.Category) error
}

var _ CategoryRepo = (*categoryRepo)(nil)

type categoryRepo struct {
	DB Queryable
}

func NewCategoryRepo(queryable Queryable) *categoryRepo {
	return &categoryRepo{queryable}
}

func (c *categoryRepo) Insert(category model.Category) (id int64, err error) {
	query := `INSERT INTO category(categoryName) VALUES(?)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := c.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, category.CategoryName)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// InsertMany 插入多个
func (c *categoryRepo) InsertMany(list []model.Category) error {
	query := `INSERT INTO category(categoryName) VALUES `

	var inserts []string
	var params []interface{}
	for _, v := range list {
		inserts = append(inserts, "(?)")
		params = append(params, v.CategoryName)
	}
	queryVals := strings.Join(inserts, ",")
	query += queryVals
	log.Println("insert sql: ", query)
	log.Println("insert params: ", params)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stmt, err := c.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("InsertMany() 准备SQL时错误: ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		log.Println("InsertMany() 执行SQL时错误: ", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("InsertMany() 获取影响行时错误: ", err)
		return err
	}

	log.Printf("category表创建 %d 条数据", rows)

	return nil
}
