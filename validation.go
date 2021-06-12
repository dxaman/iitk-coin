package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkUser(roll string) bool{
	var flag = false
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS student (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT)")
	checkErr(err)
	statement.Exec()
	rows, _ := database.Query("SELECT rollno FROM student")
	var rollno string
	for rows.Next() {
		rows.Scan(&rollno)
		if rollno  == roll {
			 flag = true
		}
	}
	return flag
}

func checkPassword(roll string, pass string) bool{
	var flag = false
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS student (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT)")
	checkErr(err)
	statement.Exec()
	rows, _ := database.Query("SELECT rollno,password FROM student")
	var password,rollno string
	for rows.Next() {
		rows.Scan(&rollno, &password)
		if rollno==roll && comparePasswords(password,[]byte(pass)){
			flag = true
		}
	}
	return flag
}
