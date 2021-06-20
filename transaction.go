package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"sync"
)
type Transactions struct {
	Coins int `json:"coins"`
	To string `json:"to"`
}
func Balance(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("token")
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
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var availBal  = fetchBal(claims.Rollno)
	if availBal != -1{
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(fmt.Sprintf("You have %s coins!", strconv.Itoa(availBal))))
		return
	}
	tx.Rollback()
	w.Write([]byte(fmt.Sprintf("NO DATA AVAILABLE!")))
	return
}

func Transfer(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("token")
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
	var transactions Transactions
	errr := json.NewDecoder(r.Body).Decode(&transactions)
	if errr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if transactions.Coins < 0{
		w.Write([]byte("Unsupported Amount!\n"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	if waitList("fetch",claims.Rollno,0,tx,&wg)>=transactions.Coins && transactions.To!=claims.Rollno{
		wg.Wait()
		wg.Add(2)
		if waitList("update",transactions.To, transactions.Coins,tx,&wg) ==1 && waitList("update",claims.Rollno,-transactions.Coins,tx,&wg)==1{
			wg.Wait()
			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}
			//sleep()
			wg.Add(1)
			//var availBal = waitList("fetch",transactions.To,0,tx,&wg)
			var ownerBal = waitList("fetch",claims.Rollno,0,tx,&wg)
			wg.Wait()
			w.Write([]byte(fmt.Sprintf("You have %s coins left!",strconv.Itoa(ownerBal))))
			return
		}
		w.Write([]byte(fmt.Sprintf("Unexpected Error Occured")))
		tx.Rollback()
		return
	}
	w.Write([]byte(fmt.Sprintf("Insufficient Balance!")))
	tx.Rollback()
	return

}

func Award(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("token")
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
	var transactions Transactions
	errr := json.NewDecoder(r.Body).Decode(&transactions)
	if errr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	database, _ := sql.Open("sqlite3", "./data_dxaman_0.db")
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	if waitList("update",transactions.To, transactions.Coins,tx,&wg) ==1{
		wg.Wait()
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(fmt.Sprintf("Awardee has been awarded %s coins!", strconv.Itoa(transactions.Coins))))
		return
	}
	tx.Rollback()
	w.Write([]byte(fmt.Sprintf("Unexpected Error Occured")))
	return

}