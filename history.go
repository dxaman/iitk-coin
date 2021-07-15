package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func addToHistory(query string, cRoll string , tTo string , tCoins int, status string, reqID int ){
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY,requestID int,mode TEXT,status TEXT, rollnoFrom TEXT,rollnoTo TEXT, coinsTransferred INT,date TEXT, time TEXT)")
	statement.Exec()
	t := time.Now()
	statement, _ = database.Prepare("INSERT INTO history (requestID,mode,status,rollnoFrom,rollnoTo,coinsTransferred,date,time)  VALUES (?,?,?,?,?,?,?,?)")
	statement.Exec(reqID, query,status,cRoll,tTo,tCoins,t.Format("02-01-2006"),t.Format("15:04:05"))
	return
}
func History(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false"{
		fetchHistory(authUser,w,r)
	}
}
func adminHistory(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser) {
		var admin Admin
		errr := json.NewDecoder(r.Body).Decode(&admin)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fetchHistory(admin.Rollno,w,r)
	}
}
func fetchHistory(cRoll string,w http.ResponseWriter, r *http.Request){
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY,requestID int,mode TEXT,status TEXT, rollnoFrom TEXT,rollnoTo TEXT, coinsTransferred INT,date TEXT, time TEXT)")
	checkErr(err)
	statement.Exec()
	var requestID,coins int
	var rollnoFrom,status, rollnoTo,mode,date,time string
	rows, _ := database.Query("SELECT requestID,mode,status,rollnoFrom,rollnoTo,coinsTransferred,date,time FROM history")
	for rows.Next(){
		rows.Scan(&requestID,&mode,&status, &rollnoFrom, &rollnoTo,&coins,&date,&time)
		if (rollnoFrom == cRoll || rollnoTo == cRoll){
			w.Write([]byte(fmt.Sprintf("Request ID: "+ strconv.Itoa(requestID) + "\nMode: "+ mode + "\nStatus: "+status+" \nFrom: "+ rollnoFrom+" \nTo: "+ rollnoTo+" \nCoins: " +strconv.Itoa(coins)+"\nDate: "+ date + " \nTime: "+ time+"\n")))
		}
	}
	defer database.Close()
}
func updateHis(reqid int , status string){
	    database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, err := database.Prepare("UPDATE history SET status = ?  WHERE requestID = ?")
		if err != nil {
			return
		}
		_,err = statement.Exec(status, reqid)
		if err != nil {
			return
		}
		return
}
