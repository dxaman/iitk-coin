package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)
type Admin struct {
	Rollno string `json:"rollno"`
}
type Storelist struct{
	Item string `json:"item"`
	ItemId string `json:"itemid"`
	Price int `json:"price"`
}
type Request struct{
	Status string `json:"status"`
	ReqId int `json:"reqid"`

}
func adminHome(w http.ResponseWriter, r *http.Request) {
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser){
		w.Write([]byte(fmt.Sprintf("Hello Admin, %s", authUser)))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func Store(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser){
		var storelist Storelist
		errr := json.NewDecoder(r.Body).Decode(&storelist)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS store (id INTEGER PRIMARY KEY, item TEXT,itemid TEXT , price INT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO store (item,itemid,price)  VALUES (?,?,?)")
		statement.Exec(storelist.Item,storelist.ItemId, storelist.Price)
		fetchStore(w,r)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func RedeemExe(w http.ResponseWriter, r *http.Request,reqid int){
	var itemid = fetchItem(reqid)
	var rollno = fetchRollno(reqid)
	var price = fetchPrice(itemid)
	if global("redeem", rollno, "NULL", price, w) != -1 {
		updateHis(reqid,"APPROVED")
		deleteRedeem(reqid)
		w.Write([]byte(fmt.Sprintf("Transaction Approved!")))
		return
	}
	w.Write([]byte(fmt.Sprintf("NOT ELIGIBLE TO REDEEM!")))
	return

}

func makeAdmin(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser){
		var admin Admin
		errr := json.NewDecoder(r.Body).Decode(&admin)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY, rollno TEXT,coins INT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO admins (rollno)  VALUES (?)")
		statement.Exec(admin.Rollno)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var authUser = checkAuth(w, r)
	if authUser != "false" && checkAdmin(authUser) {
		var admin Admin
		errr := json.NewDecoder(r.Body).Decode(&admin)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("DELETE FROM college WHERE rollno = ?")
		statement.Exec(admin.Rollno)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	var authUser = checkAuth(w, r)
	if authUser != "false" && checkAdmin(authUser) {
		var admin Admin
		errr := json.NewDecoder(r.Body).Decode(&admin)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("DELETE FROM admins WHERE rollno = ?")
		statement.Exec(admin.Rollno)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func deleteRedeem(reqid int) {
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("DELETE FROM redeem WHERE requestID = ?")
		statement.Exec(reqid)
		return
}
func Freeze(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w,r)
	if authUser!="false" && checkAdmin(authUser){
		var admin Admin
		errr := json.NewDecoder(r.Body).Decode(&admin)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS frozen (id INTEGER PRIMARY KEY, rollno TEXT,coins INT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO frozen (rollno)  VALUES (?)")
		statement.Exec(admin.Rollno)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	return
}
func fetchWal() int {
	var flag = -1
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT rollno,coins FROM admins")
	var rollno string
	var coins int
	for rows.Next() {
		err := rows.Scan(&rollno,&coins)
		if err != nil {
			return -1
		}
		if rollno  == "wallet" {
			flag = coins
		}
	}
	return flag
}
func updateWal(amt int,tx *sql.Tx) int{
	var flag = -1
	var availWal = fetchWal()
	if availWal != -1{
		statement, err := tx.Prepare("UPDATE admins SET coins = ?  WHERE rollno = ?")
		if err != nil {
			return -1
		}
		_,err = statement.Exec(availWal+amt, "wallet")
		if err != nil {
			return -1
		}
		flag =  1
	}
	return flag
}
func fetchReqID() int{
	var flag = -1
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT rollno,coins FROM admins")
	var rollno string
	var coins int
	for rows.Next() {
		err := rows.Scan(&rollno,&coins)
		if err != nil {
			return -1
		}
		if rollno  == "requestID" {
			flag = coins
		}
	}
	if incrementID(flag,database) == -1{
		flag = -1
	}
	return flag
}
func incrementID(amt int,db *sql.DB) int{
	var flag = -1
	var availWal = fetchWal()
	if availWal != -1{
		statement, err := db.Prepare("UPDATE admins SET coins = ?  WHERE rollno = ?")
		if err != nil {
			return -1
		}
		_,err = statement.Exec(amt+1, "requestID")
		if err != nil {
			return -1
		}
		flag =  1
	}
	return flag
}

func fetchRedeem(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w, r)
	if authUser != "false" && checkAdmin(authUser) {
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS redeem (id INTEGER PRIMARY KEY,requestID INT, rollno TEXT, itemid TEXT, status TEXT)")
		checkErr(err)
		statement.Exec()
		var  requestID int
		var itemid, rollno,status string
		rows, _ := database.Query("SELECT requestID,rollno,itemid,status FROM redeem")
		for rows.Next() {
			rows.Scan(&requestID,&rollno , &itemid, &status)
			w.Write([]byte(fmt.Sprintf("Request ID: " + strconv.Itoa(requestID) + "\nRoll No: " + rollno + "\nItem ID: " + itemid + "\nStatus:" + status + "\n")))
		}
		defer database.Close()
	}
}
func fetchItem(reqid int)string{
	var flag = "-1"
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT requestID,itemid FROM redeem")
	var itemid string
	var requestID int
	for rows.Next() {
		err := rows.Scan(&requestID,&itemid)
		if err != nil {
			return "-1"
		}
		if reqid  == requestID {
			flag = itemid
		}
	}
	return flag
}
func fetchRollno(reqid int)string{
	var flag = "-1"
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT requestID,rollno FROM redeem")
	var rollno string
	var requestID int
	for rows.Next() {
		err := rows.Scan(&requestID,&rollno)
		if err != nil {
			return "-1"
		}
		if reqid  == requestID {
			flag = rollno
		}
	}
	return flag
}
func fetchPrice(item string)int{
	var flag = -1
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	rows, _ := database.Query("SELECT itemid,price FROM store")
	var itemid string
	var price int
	for rows.Next() {
		err := rows.Scan(&itemid,&price)
		if err != nil {
			return -1
		}
		if itemid  == item {
			flag = price
		}
	}
	return flag
}

func Approve(w http.ResponseWriter, r *http.Request){
	var authUser = checkAuth(w, r)
	if authUser != "false" && checkAdmin(authUser) {
		var request Request
		errr := json.NewDecoder(r.Body).Decode(&request)
		if errr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		takeAction(request.ReqId,request.Status,w,r)
	}

}
func takeAction(reqid int, status string,w http.ResponseWriter, r *http.Request){
	if status == "APPROVED"{
		RedeemExe(w,r,reqid)
	} else{
		updateHis(reqid,"REJECTED")
		deleteRedeem(reqid)
		w.Write([]byte(fmt.Sprintf("Request Rejected!")))
	}
}
