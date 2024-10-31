package sqlhelper

import (
	"fmt"
	"reflect"
	"strings"
)

// Join将多条sql通过';'拼接成一条sql。
func JoinSQLs(sqls ...string) string {
	for i, sql := range sqls {
		sqls[i] = strings.TrimSuffix(strings.TrimSpace(sql), ";")
	}
	return strings.Join(sqls, ";\n")
}

func JoinArgs(args ...[]any) []any {
	if len(args) == 0 {
		return nil
	}
	n := 0
	for _, as := range args {
		n += len(as)
	}
	if n == 0 {
		return nil
	}
	all := make([]any, 0, n)
	for _, as := range args {
		all = append(all, as...)
	}
	return all
}

// ToAnySlice将任意类型的数组转为[]any。
func ToAnySlice[T any](ls []T) []any {
	results := make([]any, len(ls))
	for i := range results {
		results[i] = ls
	}
	return results
}

// InMarks用于生成sql中的in语句。
// 返回len(ls)个?。
// 用法：
//
//	args := []uint64{1, 2, 3}
//	query := fmt.Sprintf("select * from db.table where id in (%v)", InMarks(len(args)))
//	var results []*Struct
//	Query(db, query, ToAnySlice(args), &results)
func InMarks(n int) string {
	marks := make([]string, n)
	for i := range marks {
		marks[i] = "?"
	}
	return strings.Join(marks, ", ")
}

// InString用于生成sql中的in语句。
// InString返回的字符串不能用于替换sql中的'?'，
// 只能直接用于sql字符串中。
// 如：
//
//	fmt.Sprintf("select * from db.table where id in (%s)", InString([]uint64{1, 2, 3}))。
func InString[T any](ls []T) string {
	if reflect.TypeOf(ls).Elem().Kind() != reflect.String {
		return strings.Join(fmtStrs(ls), ", ")
	}
	ss := fmtStrs(ls)
	for i, s := range ss {
		ss[i] = Escape(s)
	}
	return "'" + strings.Join(ss, "', '") + "'"
}

func fmtStrs[T any](ls []T) []string {
	if len(ls) == 0 {
		return nil
	}
	ss := make([]string, len(ls))
	for i := range ls {
		ss[i] = fmt.Sprint(ls[i])
	}
	return ss
}
