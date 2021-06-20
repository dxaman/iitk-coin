package main

import (
	"database/sql"
	"encoding/json"
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

	//Un-comment the below code only if want to see all the entries in the database with their hashed password.

	/*database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS college (id INTEGER PRIMARY KEY, rollno TEXT, fullname TEXT, password TEXT,coins INT)")
	checkErr(err)
	statement.Exec()
	var id,coins int
	var rollno, fullname,password string
	rows, _ := database.Query("SELECT id, rollno, fullname,password,coins FROM college")
	for rows.Next(){
		rows.Scan(&id, &rollno, &fullname,&password,&coins)
		fmt.Println(strconv.Itoa(id) + ": "+ fullname + " - "+ rollno+" - "+ password+" -" +strconv.Itoa(coins))
	}
	defer database.Close()*/

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
	checkAuth(w,r)

/*	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
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
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Rollno)))*/

}
