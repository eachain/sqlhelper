package sqlhelper

import (
	"fmt"
	"reflect"
	"strings"
)

// GenInsertSQL生成一条插入一行记录的insert语句。
func GenInsertSQL(table string, src any) (string, []any) {
	if table != "" && table[0] != '`' {
		table = "`" + table
	}
	if table != "" && table[len(table)-1] != '`' {
		table += "`"
	}

	val := reflect.Indirect(reflect.ValueOf(src))
	typ := val.Type()
	n := typ.NumField()
	columns := make([]string, 0, n)
	args := make([]any, 0, n)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		fs := strings.Split(tag, ",")
		if len(fs) > 1 && fs[1] == "readonly" {
			continue
		}
		columns = append(columns, fs[0])
		field := val.Field(i)
		args = append(args, field.Interface())
	}

	return fmt.Sprintf("insert into %v (`%v`) values (%v)",
		table, strings.Join(columns, "`, `"), InMarks(len(args))), args
}

// GenMultiInsertSQL生成一条插入多行记录的insert语句。
func GenMultiInsertSQL[T any](table string, srcs []T) (string, []any) {
	if len(srcs) == 0 {
		return "", nil
	}

	if table != "" && table[0] != '`' {
		table = "`" + table
	}
	if table != "" && table[len(table)-1] != '`' {
		table += "`"
	}

	columns := genInsertFields(srcs[0])
	var names []string
	var args []any
	for _, src := range srcs {
		as := genInsertValues(src)
		names = append(names, "("+InMarks(len(as))+")")
		args = append(args, as...)
	}

	return fmt.Sprintf("insert into %v (`%v`) values %v",
		table, strings.Join(columns, "`, `"), strings.Join(names, ", ")), args
}

func genInsertFields(src any) []string {
	val := reflect.Indirect(reflect.ValueOf(src))
	typ := val.Type()
	n := typ.NumField()
	columns := make([]string, 0, n)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		fs := strings.Split(tag, ",")
		if len(fs) > 1 && fs[1] == "readonly" {
			continue
		}
		columns = append(columns, fs[0])
	}

	return columns
}

func genInsertValues(src any) []any {
	val := reflect.Indirect(reflect.ValueOf(src))
	typ := val.Type()
	n := typ.NumField()
	args := make([]any, 0, n)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		fs := strings.Split(tag, ",")
		if len(fs) > 1 && fs[1] == "readonly" {
			continue
		}
		args = append(args, val.Field(i).Interface())
	}

	return args
}
