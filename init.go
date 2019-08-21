package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// BIRDTODO: 程序一开始从 coverHub中获取存在些什么文件
func doInit() {
	fmt.Print(`初始化数据库意味着将删除现有的数据库的一切（如果现在有数据库），输入Enter以继续...`)
	var answer string
	var pass = false
	_, err := fmt.Scanln(&answer)
	if err != nil && err.Error() == "unexpected newline" {
		pass = true
	} else if err != nil {
		panic(err)
	}

	if pass {
		fmt.Println(`>> 开始初始化数据库`)
		os.Remove("./coverHub.db")
		db, err := sql.Open("sqlite3", "./coverHub.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStm := `
		create table url2Path ( url text not null primary key, path text );
		delete from url2Path;
		`
		_, err = db.Exec(sqlStm)
		if err != nil {
			panic(err)
		}
		fmt.Println(`>> 数据库初始化完毕`)
	} else {
		fmt.Println(`>> 跳过 数据库初始化`)
	}

	fmt.Print(`初始化图片仓库意味着将删除现有图片仓库的一切（如果现在有图片仓库），输入Enter以继续...`)
	_, err = fmt.Scanln(&answer)
	pass = false
	if err != nil && err.Error() == "unexpected newline" {
		pass = true
	} else if err != nil {
		panic(err)
	}

	if pass {
		fmt.Println(`>> 开始初始化图片仓库`)
		err = os.RemoveAll("./coverFiles/")
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll("./coverFiles/temp", os.ModePerm)
		if err != nil {
			panic(err)
		}
		fmt.Println(`>> 图片仓库初始化完毕`)
	} else {
		fmt.Println(`>> 程序结束`)
	}
}
