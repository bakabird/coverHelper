package main

import (
	"database/sql"
	"flag"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// useage of go-sqlite3 can be known on https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go

const coverFileHub string = `./coverFiles/`
const coverFileTempHub string = `./coverFiles/temp`
const dbPath string = "./coverHub.db"

var serverTimes = 0
var mu_db sync.Mutex

var DOINIT = flag.Bool("i", false, "重新初始化数据库和图片仓库")
var SHOW_VERSION = flag.Bool("v", false, "打印目前的版本号")
var STATIC_COVER_HUB = flag.String("huburl", "http://127.0.0.1:8360/static/", "图片仓库对应的静态路径，带URL")

// BIRDTODO: 避免同时下载同一个文件。 <-- 用一个MAP来记录工人准备去下载的
func main() {
	flag.Parse()
	if *SHOW_VERSION {
		fmt.Println("version: 0.1.2")
	} else if *DOINIT {
		doInit()
	} else {
		fmt.Println("huburl", *STATIC_COVER_HUB)

		// BIRDTODO: 在init程序中确保图片仓库存在
		db, err := sql.Open("sqlite3", dbPath)
		fmt.Println("开启数据库连接")
		defer db.Close()
		defer fmt.Println("数据库")
		if err != nil {
			fmt.Println("在启动前请确保数据库存在")
			return
		}
		startListen(db)
	}
}
