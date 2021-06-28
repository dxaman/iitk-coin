package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)
type Admin struct {
	Rollno string `json:"rollno"`
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
