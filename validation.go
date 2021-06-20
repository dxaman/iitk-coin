package main

import (
	//"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"time"

	//"log"
	"net/http"
)
var mutex sync.Mutex
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func sleep(){
	time.Sleep(4 * time.Second)
}

func checkUser(roll string) bool{
	var flag = false
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS college (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT)")
	checkErr(err)
	statement.Exec()
	rows, _ := database.Query("SELECT rollno FROM college")
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
	rows, err := database.Query("SELECT rollno,password FROM college")
	checkErr(err)
	var password,rollno string
	for rows.Next() {
		rows.Scan(&rollno, &password)
		if rollno==roll && comparePasswords(password,[]byte(pass)){
			flag = true
		}
	}
	return flag
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool{
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return false
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Rollno)))
	return true

}
func waitList(query string,roll string, amt int , tx *sql.Tx,wg *sync.WaitGroup)int{
	mutex.Lock()
	if query=="fetch"{
		var val = fetchBal(roll)
		mutex.Unlock()
		wg.Done()
		return val
	}
	if query=="update"{
		var availCoins = fetchBal(roll)
		//sleep()
		var flag = updateBal(roll,amt,availCoins,tx)
		mutex.Unlock()
		wg.Done()
		return flag
	}
	return -1
}

func fetchBal(roll string) int {
	var flag = -1
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT rollno,coins FROM college")
	var rollno string
	var coins int
	for rows.Next() {
		err := rows.Scan(&rollno,&coins)
		if err != nil {
			return -1
		}
		if rollno  == roll {
			flag = coins
		}
	}
	return flag
}
func updateBal(roll string,amt int,availCoins int,tx *sql.Tx) int{
	var flag = -1
	//var availCoins = fetchBal(roll)
	if availCoins != -1{
		statement, err := tx.Prepare("UPDATE college SET coins = ?  WHERE rollno = ?")
		if err != nil {
			return -1
		}
		_,err = statement.Exec(availCoins+amt, roll)
		if err != nil {
			return -1
		}
		flag =  1
	}
	return flag
}
