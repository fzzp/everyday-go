package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// execTx 执行事务
func execTx(ctx context.Context, qb Queryable, fn func(*Repository) error) error {
	// qb => sql.DB/sql.Tx
	db, ok := qb.(*sql.DB)
	if !ok {
		return errors.New("qb Queryable 不是 *sql.DB 指针")
	}

	// 开启事务
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repo := NewRepository(tx)
	if err = fn(repo); err != nil {
		// 执行失败，回滚
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("执行事务失败: %s, rb err: %s", err, rbErr)
		}
		return err
	}

	// 执行成功 提交
	return tx.Commit()
}
