
package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/award", Award)
	http.HandleFunc("/transfer", Transfer)
	http.HandleFunc("/balance", Balance)
	http.HandleFunc("/redeem", Redeem)
	http.HandleFunc("/admin", adminHome)
	http.HandleFunc("/admin/makeAdmin", makeAdmin)
	http.HandleFunc("/admin/deleteUser", deleteUser)
	http.HandleFunc("/admin/deleteAdmin", deleteAdmin)
	http.HandleFunc("/admin/storeList", Store)
	http.HandleFunc("/admin/redeemReq", fetchRedeem)
	http.HandleFunc("/admin/approval", Approve)
	http.HandleFunc("/admin/freeze", Freeze)
	http.HandleFunc("/history", History)
	http.HandleFunc("/admin/history", adminHistory)
	http.HandleFunc("/database", fetchDatabase)
	http.HandleFunc("/store", fetchStore)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
