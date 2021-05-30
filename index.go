package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)


func main() {
	database, _ := sql.Open("sqlite3","./data_dxaman_1.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS student (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO student (rollno, fullname)  VALUES (?,?)")
	statement.Exec("190558","Shruti Nisal")
	var id int
	var rollno, fullname string
	rows, _ := database.Query("SELECT id, rollno, fullname FROM student")
	for rows.Next(){
		rows.Scan(&id, &rollno, &fullname)
		fmt.Println(strconv.Itoa(id) + ": "+ fullname + " - "+ rollno)
	}


}
