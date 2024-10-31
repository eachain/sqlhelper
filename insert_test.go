package sqlhelper

import (
	"testing"
	"time"
)

type person struct {
	Id       uint64    `db:"id,readonly"`
	Name     string    `db:"name"`
	Age      int       `db:"age"`
	CreateAt time.Time `db:"create_at,readonly"`
}

func TestGenInsertSQL(t *testing.T) {
	me := &person{
		Id:       1,
		Name:     "eachain",
		Age:      1024,
		CreateAt: time.Now(),
	}
	sql, args := GenInsertSQL("person", me)

	query := "insert into `person` (`name`, `age`) values (?, ?)"
	if sql != query {
		t.Fatalf("insert sql: %v", sql)
	}

	if len(args) != 2 {
		t.Fatalf("args: %+v", args)
	}

	if args[0] != "eachain" {
		t.Fatalf("first arg: %+v", args[0])
	}

	if args[1] != 1024 {
		t.Fatalf("second arg: %+v", args[1])
	}
}

func TestGenMultiInsertSQL(t *testing.T) {
	ps := []*person{{
		Id:       1,
		Name:     "eachain",
		Age:      1024,
		CreateAt: time.Now(),
	}, {
		Id:       1,
		Name:     "foolish",
		Age:      999,
		CreateAt: time.Now(),
	}}
	sql, args := GenMultiInsertSQL("person", ps)

	query := "insert into `person` (`name`, `age`) values (?, ?), (?, ?)"
	if sql != query {
		t.Fatalf("insert sql: %v", sql)
	}

	if len(args) != 4 {
		t.Fatalf("args: %+v", args)
	}
	if args[0] != "eachain" {
		t.Fatalf("first arg: %+v", args[0])
	}
	if args[1] != 1024 {
		t.Fatalf("second arg: %+v", args[1])
	}
	if args[2] != "foolish" {
		t.Fatalf("third arg: %+v", args[0])
	}
	if args[3] != 999 {
		t.Fatalf("fourth arg: %+v", args[1])
	}
}
