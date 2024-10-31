# sqlhelper

sqlhelper实现了一组简化标准库`database/sql`的方法。可以理解为轻量级orm。

## Features

- `Query`结果自动映射到结构体；
	- 标签：`db`，如果某些字段在生成insert语句时不需要，可以在`db`的值里面加上`readonly`。
- 支持解析多条sql结果；
- 根据结构体tag生成insert语句；

## 示例

```go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eachain/sqlhelper"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	// `id`字段自增，Query时返回具体值，生成insert语句时忽略
	Id         uint64    `db:"id,readonly" json:"id"`
	Name       string    `db:"name" json:"name"`
	Gender     int       `db:"gender" json:"gender"`
	Email      string    `db:"email" json:"email"`
	CreateTime time.Time `db:"create_time,readonly" json:"create_time"`
}

func main() {
	dsn := (&mysql.Config{
		User:                 "user",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		Loc:                  time.Local,
		Timeout:              3 * time.Second,
		AllowNativePasswords: true,
		MultiStatements:      true, // 需要将该值设为true才允许执行多条sql
		InterpolateParams:    true, // 需要将该值设为true才允许多条sql带'?'参数
		ParseTime:            true,
	}).FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query1 := "select count(1) from user where gender = ?"
	args1 := []any{0}
	var count int

	query2 := "select * from user where id in (?, ?)"
	args2 := []any{1, 2}
	var users []*User

	// 执行多条sql，并将结果解析到对应结构中
	err = sqlhelper.Query(
		db,
		sqlhelper.JoinSQLs(query1, query2),
		sqlhelper.JoinArgs(args1, args2),
		&count, &users)
	if err != nil {
		panic(err)
	}

	fmt.Printf("gender users count: %v\n", count)

	js, _ := json.Marshal(users)
	fmt.Printf("users: %s\n", js)

	fmt.Println(sqlhelper.GenMultiInsertSQL("user", users))
	// 生成的insert语句中，将忽略标有`readonly`的字段：
	// insert into `user` (`name`, `gender`, `email`) values (?, ?, ?), (?, ?, ?)
}
```
