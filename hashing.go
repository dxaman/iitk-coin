package main

import (
	"golang.org/x/crypto/bcrypt"
)
func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	checkErr(err)
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		//log.Println(err)
		return false
	}

	return true
}