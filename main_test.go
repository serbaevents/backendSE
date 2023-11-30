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
	hasil, err := watoken.Encode("cobainah", privateKey)
	fmt.Println(hasil, err)
}

func TestInsertUser(t *testing.T) {
	mconn := GetConnectionMongo("MONGOSTRING", "well")
	var userdata User
	userdata.Username = "bdSE"
	userdata.Password = "cobain"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestInsertAdmin(t *testing.T) {
	mconn := GetConnectionMongo("MONGOSTRING", "User")
	var userdata Admin
	userdata.Username = "admin"
	userdata.Password = "cobain"
	userdata.Role = "admin"

	nama := InsertAdmin(mconn, "admin", userdata)
	fmt.Println(nama)
}
