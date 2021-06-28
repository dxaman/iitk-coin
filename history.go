package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func addToHistory(query string, cRoll string , tTo string , tCoins int ){
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY,mode TEXT, rollnoFrom TEXT,rollnoTo TEXT, coinsTransferred INT,date TEXT, time TEXT)")
	statement.Exec()
	t := time.Now()
	statement, _ = database.Prepare("INSERT INTO history (mode,rollnoFrom,rollnoTo,coinsTransferred,date,time)  VALUES (?,?,?,?,?,?)")
	statement.Exec(query,cRoll,tTo,tCoins,t.Format("02-01-2006"),t.Format("15:04:05"))
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
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT,coins INT)")
	checkErr(err)
	statement.Exec()
	var id,coins int
	var rollnoFrom, rollnoTo,mode,date,time string
	rows, _ := database.Query("SELECT id,mode,rollnoFrom,rollnoTo,coinsTransferred,date,time FROM history")
	for rows.Next(){
		rows.Scan(&id,&mode, &rollnoFrom, &rollnoTo,&coins,&date,&time)
		if (rollnoFrom == cRoll || rollnoTo == cRoll){
			w.Write([]byte(fmt.Sprintf(strconv.Itoa(id) + "\nMode: "+ mode + " \nFrom: "+ rollnoFrom+" \nTo: "+ rollnoTo+" \nCoins=" +strconv.Itoa(coins)+"\nDate: "+ date + " \nTime: "+ time+"\n")))
		}
	}
	defer database.Close()
}
