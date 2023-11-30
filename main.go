package gobd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/mongo"
)

func GCHandlerFunc(Mongostring, dbname, colname string) string {
	koneksyen := GetConnectionMongo(Mongostring, dbname)
	databis := GetAllData(koneksyen, colname)

	jsoncihuy, _ := json.Marshal(databis)

	return string(jsoncihuy)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := GetConnectionMongo(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return ReturnStringStruct(Response)
}

func InsertUser(db *mongo.Database, collection string, userdata User) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Ini username : " + userdata.Username + "ini password : " + userdata.Password
}

func InsertAdmin(db *mongo.Database, collection string, userdata Admin) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Ini username : " + userdata.Username + "ini password : " + userdata.Password + "ini Role : " + userdata.Role
}

// func GCFPostCoordinate(Mongostring, dbname, colname string, r *http.Request) string {
// 	req := new(Credents)
// 	conn := GetConnectionMongo(Mongostring, dbname)
// 	resp := new(LonLatProperties)
// 	err := json.NewDecoder(r.Body).Decode(&resp)
// 	if err != nil {
// 		req.Status = strconv.Itoa(http.StatusNotFound)
// 		req.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		req.Status = strconv.Itoa(http.StatusOK)
// 		Ins := InsertDataLonlat(conn, colname,
// 			resp.Coordinates,
// 			resp.Name,
// 			resp.Volume,
// 			resp.Type)
// 		req.Message = fmt.Sprintf("%v:%v", "Berhasil Input data", Ins)
// 	}
// 	return ReturnStringStruct(req)
// }
func InsertDataBis(Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	conn := GetConnectionMongo(Mongostring, dbname)
	resp := new(Bis)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		req.Status = strconv.Itoa(http.StatusOK)
		Ins := InsertDatabisini(conn, colname,
			resp.JamGo,
			resp.JamOut,
			resp.NoKursi,
			resp.Jemputan)
		req.Message = fmt.Sprintf("%v:%v", "Berhasil Input data", Ins)
	}
	return ReturnStringStruct(req)
}
func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

// func GCFUpdate(Mongostring, dbname, colname string, r *http.Request) string {
// 	req := new(Credents)
// 	resp := new(LonLatProperties)
// 	conn := GetConnectionMongo(Mongostring, dbname)
// 	err := json.NewDecoder(r.Body).Decode(&resp)
// 	if err != nil {
// 		req.Status = strconv.Itoa(http.StatusNotFound)
// 		req.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		req.Status = strconv.Itoa(http.StatusOK)
// 		Ins := UpdateDataGeojson(conn, colname,
// 			resp.Name,
// 			resp.Volume,
// 			resp.Type)
// 		req.Message = fmt.Sprintf("%v:%v", "Berhasil Update data", Ins)
// 	}
// 	return ReturnStringStruct(req)
// }

// func GCFDeleteData(Mongostring, dbname, colname string, r *http.Request) string {
// 	req := new(Credents)
// 	resp := new(LonLatProperties)
// 	conn := GetConnectionMongo(Mongostring, dbname)
// 	err := json.NewDecoder(r.Body).Decode(&resp)
// 	if err != nil {
// 		req.Status = strconv.Itoa(http.StatusNotFound)
// 		req.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		req.Status = strconv.Itoa(http.StatusOK)
// 		delResult, delErr := DeleteDataGeojson(conn, colname, resp.Name)
// 		if delErr != nil {
// 			req.Status = strconv.Itoa(http.StatusInternalServerError)
// 			req.Message = "error deleting data: " + delErr.Error()
// 		} else {
// 			req.Message = fmt.Sprintf("Berhasil menghapus data. Jumlah data terhapus: %v", delResult.DeletedCount)
// 		}
// 	}
// 	return ReturnStringStruct(req)
// }




