package adapters

import (
	"database/sql"
	_ "database/mysql"
)

var database = []string{
	"spinner1:271312cs@/spinner?parseTime=true",
	"spinner2:271312cs@/spinner?parseTime=true",
	"spinner3:271312cs@/spinner?parseTime=true",
	"spinner4:271312cs@/spinner?parseTime=true",
	"spinner5:271312cs@/spinner?parseTime=true",
	"spinner6:271312cs@/spinner?parseTime=true",
	"spinner7:271312cs@/spinner?parseTime=true",
	"spinner8:271312cs@/spinner?parseTime=true",
}

func ConnectDB(index int) *sql.DB {
 	db, err := sql.Open("mysql", database[index - 1])
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}