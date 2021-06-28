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
var cap = 5000
var minEve = 5

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
		if fetchEve(authUser)>=minEve{
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
			if global("transfer", authUser, transactions.To, transactions.Coins, w) ==1{
				addToHistory("transfer",authUser,transactions.To,transactions.Coins)
			}
			return
		}
		w.Write([]byte(fmt.Sprintf("NOT ELIGIBLE TO TRANSFER!")))
	}
}

func Award(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser){
		var transactions Transactions
		errr := json.NewDecoder(r.Body).Decode(&transactions)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if global("award", authUser, transactions.To, transactions.Coins, w)==1{
			addToHistory("reward","NULL",transactions.To,transactions.Coins)
			eveInc(transactions.To)
			return
		}
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func Redeem(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" {
		if fetchEve(authUser)>=minEve {
			var transactions Transactions
			errr := json.NewDecoder(r.Body).Decode(&transactions)
			if errr != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if global("redeem", authUser, transactions.To, transactions.Coins, w) == 1 {
				addToHistory("redeem", authUser, "NULL", transactions.Coins)
			}
			return
		}
		w.Write([]byte(fmt.Sprintf("NOT ELIGIBLE TO REDEEM!")))
	}

}
func fetchDatabase(w http.ResponseWriter, r *http.Request){
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS college (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT,coins INT,events INT)")
	checkErr(err)
	statement.Exec()
	var id,coins,events int
	var rollno, fullname,password string
	rows, _ := database.Query("SELECT id, rollno, fullname,password,coins,events FROM college")
	for rows.Next(){
		rows.Scan(&id, &rollno, &fullname,&password,&coins,&events)
		w.Write([]byte(fmt.Sprintf(strconv.Itoa(id) + "\nName: "+ fullname + " \nRoll Number: "+ rollno+" \nHashed Password: "+ password+" \nCoins=" +strconv.Itoa(coins)+ "\nEvents=" +strconv.Itoa(events)+"\n")))
	}
	defer database.Close()
}
