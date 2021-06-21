package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)


var jwtKey = []byte("secret_key")

type Credentials struct {
	Rollno   string `json:"rollno"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

type Claims struct {
	Rollno string `json:"rollno"`
	jwt.StandardClaims
}


func Signup(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	if !checkUser(credentials.Rollno) {
		database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS college (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT,coins INT)")
		statement.Exec()
		statement, err = database.Prepare("INSERT INTO college (rollno, fullname, password, coins)  VALUES (?,?,?,?)")
		checkErr(err)
		statement.Exec(credentials.Rollno, credentials.Fullname, hashAndSalt([]byte(credentials.Password)),0)
		w.Write([]byte("User Successfully Registered\n"))
	} else{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User Already Exists!\n"))
	}
	return

}
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !checkPassword(credentials.Rollno, credentials.Password){
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Username and Password not matched!\n"))
		return
	}
	//Generate token and set cookies
	expirationTime := time.Now().Add(time.Minute * 20)

	claims := &Claims{
		Rollno: credentials.Rollno,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c := http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
	}
	http.SetCookie(w,&c)
	w.Write([]byte("Successfully Logged In!\n"))
}
func Logout(w http.ResponseWriter, r *http.Request) {

	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
		//MaxAge<0 deletes the cookie
	http.SetCookie(w, &c)
	w.Write([]byte("Logged out!\n"))
}

func Home(w http.ResponseWriter, r *http.Request) {
	var authUser = checkAuth(w,r)
	if authUser!="false"{
		w.Write([]byte(fmt.Sprintf("Hello, %s", authUser)))
	}
}
