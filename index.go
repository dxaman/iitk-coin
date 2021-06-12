
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
