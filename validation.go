package main

import (
	"context"
	//"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
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

func checkAuth(w http.ResponseWriter, r *http.Request) string{
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return "false"
		}
		w.WriteHeader(http.StatusBadRequest)
		return "false"
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
			return "false"
		}
		w.WriteHeader(http.StatusBadRequest)
		return "false"
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return "false"
	}

	return claims.Rollno

}
func global(query string,cRollno string,tTo string, tCoins int,w http.ResponseWriter){
	mutex.Lock()
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if query=="transfer"{
		if waitList("fetch",cRollno,0,tx)>=tCoins && tTo!=cRollno{
			if waitList("update",tTo, tCoins,tx) ==1 && waitList("update",cRollno,-tCoins,tx)==1{
				//sleep()
				err = tx.Commit()
				if err != nil {
					log.Fatal(err)
				}
				var ownerBal = waitList("fetch",cRollno,0,tx)
				w.Write([]byte(fmt.Sprintf("You have %s coins left!\n", strconv.Itoa(ownerBal))))
				mutex.Unlock()
				return
			}
			w.Write([]byte(fmt.Sprintf("Unexpected Error Occured")))
			tx.Rollback()
			mutex.Unlock()
			return
		}
		w.Write([]byte(fmt.Sprintf("Insufficient Balance!")))
		tx.Rollback()
		mutex.Unlock()
		return
	}
	if query=="award"{
		if waitList("update",tTo, tCoins,tx) ==1{
			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}
			w.Write([]byte(fmt.Sprintf("Awardee has been awarded %s coins!\n", strconv.Itoa(tCoins))))
			mutex.Unlock()
			return
		}
		tx.Rollback()
		w.Write([]byte(fmt.Sprintf("Unexpected Error Occured")))
		mutex.Unlock()
		return

	}
	mutex.Unlock()
	return
}
func waitList(query string,roll string, amt int , tx *sql.Tx)int{
	if query=="fetch"{
		var val = fetchBal(roll)
		return val
	}
	if query=="update"{
		var availCoins = fetchBal(roll)
		var flag = updateBal(roll,amt,availCoins,tx)
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
