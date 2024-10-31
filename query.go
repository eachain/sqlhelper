package sqlhelper

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
)

type Querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// Query可用于(*sql.DB).Query和(*sql.Tx).Query。
// 参数query可以是单条或多条select sql通过';'拼接得来；
// dsts为解析结果，每条sql对应一个dst。
// dst支持struct结构和基础类型。
// eg. Query(db, sqls, args, &structure, &intSlice)。
func Query(db Querier, query string, args []any, dsts ...any) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	return scanRows(rows, dsts)
}

type ContextQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// QueryContext可用于(*sql.DB).QueryContext和(*sql.Tx).QueryContext。
// 参数query可以是单条或多条select sql通过';'拼接得来；
// dsts为解析结果，每条sql对应一个dst。
// dst支持struct结构和基础类型。
// eg. QueryContext(db, sqls, args, &structure, &intSlice)。
func QueryContext(db ContextQuerier, ctx context.Context, query string, args []any, dsts ...any) error {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return scanRows(rows, dsts)
}

func scanRows(rows *sql.Rows, dsts []any) error {
	defer rows.Close()

	for _, dst := range dsts {
		err := _scanRows(rows, dst)
		if err != nil {
			return err
		}
		if !rows.NextResultSet() {
			break
		}
	}

	return rows.Err()
}

func _scanRows(rows *sql.Rows, dst any) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	cols := make([]string, len(columns))
	for i := range cols {
		if idx := strings.IndexByte(columns[i], '.'); idx >= 0 {
			cols[i] = columns[i][idx+1:]
		} else {
			cols[i] = columns[i]
		}
	}

	val := reflect.ValueOf(dst).Elem()
	if val.Kind() != reflect.Slice {
		if !rows.Next() {
			return sql.ErrNoRows
		}
		return scanRow(cols, rows, val)
	}

	val.Set(val.Slice(0, 0))
	typ := val.Type().Elem()
	isPointer := false
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		isPointer = true
	}
	for rows.Next() {
		elem := reflect.New(typ)
		err := scanRow(cols, rows, elem.Elem())
		if err != nil {
			return err
		}
		if isPointer {
			val.Set(reflect.Append(val, elem))
		} else {
			val.Set(reflect.Append(val, elem.Elem()))
		}
	}
	return nil
}

func scanRow(columns []string, rows *sql.Rows, val reflect.Value) error {
	if val.Kind() != reflect.Struct {
		err := rows.Scan(val.Addr().Interface())
		if err != nil {
			return err
		}
		return nil
	}

	typ := val.Type()
	n := typ.NumField()
	idxof := make(map[string]int)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		fs := strings.Split(tag, ",")
		idxof[fs[0]] = i
	}

	addrs := make([]any, 0, len(columns))
	for _, col := range columns {
		if i, ok := idxof[col]; ok {
			addrs = append(addrs, val.Field(i).Addr().Interface())
		} else {
			addrs = append(addrs, new(any))
		}
	}
	err := rows.Scan(addrs...)
	if err != nil {
		return err
	}
	return nil
}
