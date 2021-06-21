package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)
type Transactions struct {
	Coins int `json:"coins"`
	To string `json:"to"`
}
func Balance(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" {
		var availBal = fetchBal(authUser)
		if availBal != -1 {
			w.Write([]byte(fmt.Sprintf("You have %s coins!", strconv.Itoa(availBal))))
			return
		}
		w.Write([]byte(fmt.Sprintf("NO DATA AVAILABLE!")))
		return
	}
}

func Transfer(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" {
		var transactions Transactions
		errr := json.NewDecoder(r.Body).Decode(&transactions)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if transactions.Coins < 0 {
			w.Write([]byte("Unsupported Amount!\n"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		global("transfer", authUser, transactions.To, transactions.Coins, w)
		return
	}
}

func Award(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" {
		var transactions Transactions
		errr := json.NewDecoder(r.Body).Decode(&transactions)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		global("award", authUser, transactions.To, transactions.Coins, w)
		return
	}
}
func fetchDatabase(w http.ResponseWriter, r *http.Request){
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS college (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT,coins INT)")
	checkErr(err)
	statement.Exec()
	var id,coins int
	var rollno, fullname,password string
	rows, _ := database.Query("SELECT id, rollno, fullname,password,coins FROM college")
	for rows.Next(){
		rows.Scan(&id, &rollno, &fullname,&password,&coins)
		w.Write([]byte(fmt.Sprintf(strconv.Itoa(id) + "\nName: "+ fullname + " \nRoll Number: "+ rollno+" \nHashed Password: "+ password+" \nCoins=" +strconv.Itoa(coins)+"\n")))
	}
	defer database.Close()
}
