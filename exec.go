package sqlhelper

import (
	"context"
	"database/sql"
)

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// Exec可用于(*sql.DB).Exec和(*sql.Tx).Exec。
func Exec(db Executor, query string, args ...any) (lastInsertId, rowsAffected int64, err error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return
	}
	return parseExecResult(result)
}

type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Exec可用于(*sql.DB).ExecContext和(*sql.Tx).ExecContext。
func ExecContext(db ContextExecutor, ctx context.Context, query string, args ...any) (lastInsertId, rowsAffected int64, err error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}
	return parseExecResult(result)
}

func parseExecResult(result sql.Result) (lastInsertId, rowsAffected int64, err error) {
	lastInsertId, err = result.LastInsertId()
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	return
}
