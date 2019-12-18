package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	dsName := "root:xiyou10-211@tcp(47.108.69.137:3306)/po?charset=utf8&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(3)
	db.SetConnMaxLifetime(7 * time.Hour)
	fmt.Println(db.Ping()) //成功为nil
	fmt.Println(db.Query("select now() "))

	defer db.Close()
}
