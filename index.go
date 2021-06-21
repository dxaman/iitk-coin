
package main

import (
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
	http.HandleFunc("/database", fetchDatabase)



	log.Fatal(http.ListenAndServe(":8080", nil))
}
