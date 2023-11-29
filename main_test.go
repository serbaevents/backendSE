package gobd

import (
	"fmt"
	"testing"
	"github.com/whatsauth/watoken"
)
 
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("sawcoba", privateKey)
	fmt.Println(hasil, err)
}

func TestInsertUser(t *testing.T) {
	mconn := GetConnectionMongo("MONGOSTRING", "gisaw")
	var userdata User
	userdata.Username = "sawqi"
	userdata.Password = "cobain"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}