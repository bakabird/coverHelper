package main

import (
	"database/sql"
	"fmt"
	"log"
)

func db_SavePair(db *sql.DB, url, path string) {
	fmt.Println(`DB EXEC：`, fmt.Sprintf("insert into url2Path(url, path) values(\"%s\",\"%s\")", url, path))
	_, err := db.Exec(fmt.Sprintf("insert into url2Path(url, path) values(\"%s\",\"%s\")", url, path))
	if err != nil {
		log.Fatal(err)
	}
}
func db_IsExist(db *sql.DB, url string) bool {
	rows, err := db.Query(fmt.Sprintf("select path from url2Path where url=\"%s\"", url))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("db_isExist", "由", url, "查到：", path)
		return true
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return false
}
