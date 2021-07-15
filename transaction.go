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
			var reqID = global("transfer", authUser, transactions.To, transactions.Coins, w)
			if reqID !=-1{
				addToHistory("transfer",authUser,transactions.To,transactions.Coins,"APPROVED",reqID)
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
		var reqID = global("award", authUser, transactions.To, transactions.Coins, w)
		if reqID!=-1{
			addToHistory("reward","NULL",transactions.To,transactions.Coins,"APPROVED",reqID)
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
	if authUser!="false"{
		var storelist Storelist
		errr := json.NewDecoder(r.Body).Decode(&storelist)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var reqID = global("reqID", authUser, "null", 0, w)
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS redeem (id INTEGER PRIMARY KEY,requestID INT, rollno TEXT, itemid TEXT, status TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO redeem (requestID,rollno,itemid,status)  VALUES (?,?,?,?)")
		statement.Exec(reqID,authUser,storelist.ItemId,"PENDING")
		addToHistory("redeem", authUser, "NULL", storelist.Price,"PENDING",reqID)
		w.Write([]byte(fmt.Sprintf("Request Raised! Kindly wait for Admin to take action!")))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
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
		w.Write([]byte(fmt.Sprintf(strconv.Itoa(id) + "\nName: "+ fullname + " \nRoll Number: "+ rollno+" \nHashed Password: "+ password+" \nCoins: " +strconv.Itoa(coins)+ "\nEvents: " +strconv.Itoa(events)+"\n")))
	}
	defer database.Close()
}
func fetchStore(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" {
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS store (id INTEGER PRIMARY KEY, item TEXT,itemid TEXT, price INT)")
		checkErr(err)
		statement.Exec()
		var id, price int
		var item, itemid string
		rows, _ := database.Query("SELECT id, item,itemid ,price FROM store")
		for rows.Next() {
			rows.Scan(&id, &item, &itemid, &price)
			w.Write([]byte(fmt.Sprintf(strconv.Itoa(id) + "\nItem: " + item + "\nItem ID: " + itemid + "\nPrice:" + strconv.Itoa(price) + "\n")))
		}
		defer database.Close()
	}
}
