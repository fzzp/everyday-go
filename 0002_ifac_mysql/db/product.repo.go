package db

import (
	"context"
	"database/sql"
	"everyday-go/0002_ifac_mysql/model"
	"log"
	"strings"
	"time"
)

type ProductRepo interface {
	Insert(p model.Product) (id int64, err error)
	ProductMapCategory(productID int, cateIDs []int) error
	InsertWithCategoryTx(p model.Product, list []model.Category) error
	QueryRow(id int) (model.Product, error)
}

var _ ProductRepo = (*productRepo)(nil)

type productRepo struct {
	DB Queryable
}

func NewProductRepo(queryable Queryable) *productRepo {
	return &productRepo{queryable}
}

func (repo *productRepo) Insert(p model.Product) (id int64, err error) {
	query := `INSERT INTO products(productName, productPrice) VALUES(?, ?)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.ProductName, p.ProductPrice)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func (repo *productRepo) ProductMapCategory(productID int, cateIDs []int) error {
	query := `INSERT INTO category_products(productId, categoryId) VALUES `

	var inserts []string
	var params []interface{}
	for _, id := range cateIDs {
		inserts = append(inserts, "(?, ?)")
		params = append(params, productID, id)
	}
	queryVals := strings.Join(inserts, ",")
	query += queryVals
	log.Println("insert sql: ", query)
	log.Println("insert params: ", params)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
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

// InsertWithCategoryTx 同时插入商品和分类事务
func (repo *productRepo) InsertWithCategoryTx(product model.Product, list []model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := execTx(ctx, repo.DB, func(r *Repository) error {
		newID, err := r.Products.Insert(product)
		if err != nil {
			return err
		}

		cateIDs := make([]int, 0, len(list))
		for _, category := range list {
			cateNewID, err := r.Category.Insert(category)
			if err != nil {
				return err
			}
			cateIDs = append(cateIDs, int(cateNewID))
		}

		return r.Products.ProductMapCategory(int(newID), cateIDs)

	})

	return err
}

func (repo *productRepo) QueryRow(id int) (p model.Product, err error) {
	query := `
		SELECT 
			p.id, 
			p.productName,
			p.productPrice,
			c.id,
			c.categoryName
		FROM products p 
		JOIN category_products cp ON p.id = cp.productId
		JOIN category c ON c.id = cp.categoryId 
		WHERE p.id = ?;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stmt *sql.Stmt
	stmt, err = repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.QueryContext(ctx, id)
	if err != nil {
		return
	}
	defer rows.Close()

	list := make([]model.Category, 0)
	for rows.Next() {
		var cate model.Category
		rows.Scan(&p.ID, &p.ProductName, &p.ProductPrice, &cate.ID, &cate.CategoryName)
		list = append(list, cate)
	}

	p.Categories = list

	err = rows.Err()

	return
}
