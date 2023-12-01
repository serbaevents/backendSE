// package module

// import (
// 	"encoding/json"
// 	"net/http"
// 	"os"
// 	"github.com/aiteung/atdb"
// 	"github.com/whatsauth/watoken"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // func GCFHandler(MONGOCONNSTRINGENV, dbname, collectionname string) string {
// // 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// // 	datagedung := GetAllBangunanLineString(mconn, collectionname)
// // 	return GCFReturnStruct(datagedung)
// // }

// func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	var Response Credential
// 	Response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 	} else {
// 		if IsPasswordValid(mconn, collectionname, datauser) {
// 			Response.Status = true
// 			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
// 			if err != nil {
// 				Response.Message = "Gagal Encode Token : " + err.Error()
// 			} else {
// 				Response.Message = "Selamat Datang"
// 				Response.Token = tokenstring
// 			}
// 		} else {
// 			Response.Message = "Password Salah"
// 		}
// 	}

// 	return GCFReturnStruct(Response)
// }

// func GCFReturnStruct(DataStuct any) string {
// 	jsondata, _ := json.Marshal(DataStuct)
// 	return string(jsondata)
// }

// func InsertUser(db *mongo.Database, collection string, userdata User) string {
// 	hash, _ := HashPassword(userdata.Password)
// 	userdata.Password = hash
// 	atdb.InsertOneDoc(db, collection, userdata)
// 	return "Ini username : " + userdata.Username + "ini password : " + userdata.Password
// }

package module

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/serbaevents/backendSE/model"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

// var MongoString string = os.Getenv("MONGOSTRING")

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// signup
func SignUpPengguna(db *mongo.Database, insertedDoc model.Pengguna) error {
	objectId := primitive.NewObjectID() 
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("Email sudah terdaftar")
	} 
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	} 
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "pengguna",
	}
	pengguna := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	_, err = InsertOneDoc(db, "pengguna", pengguna)
	if err != nil {
		return fmt.Errorf("kesalahan server")
	}
	return nil
}

func SignUpDriver(db *mongo.Database, insertedDoc model.Driver) error {
	objectId := primitive.NewObjectID()
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatBis == "" ||  insertedDoc.Akun.Email == "" || insertedDoc.Akun.Password == "" {
		return fmt.Errorf("dimohon untuk melengkapi data")
	} 
	if err := checkmail.ValidateFormat(insertedDoc.Akun.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	} 
	userExists, _ := GetUserFromEmail(insertedDoc.Akun.Email, db)
	if insertedDoc.Akun.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	} 
	if strings.Contains(insertedDoc.Akun.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Akun.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Akun.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"_id": objectId,
		"email": insertedDoc.Akun.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": "driver",
	}
	driver := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"platbis": insertedDoc.PlatBis,
		"akun": model.User {
			ID : objectId,
		},
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	_, err = InsertOneDoc(db, "driver", driver)
	if err != nil {
		return err
	}
	return nil
}

// login
func LogIn(db *mongo.Database, insertedDoc model.User) (user model.User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("Email tidak valid")
	} 
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return 
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

//user
func UpdateEmailUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" {
		return fmt.Errorf("Dimohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("Email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("Email sudah terdaftar")
	}
	user := bson.M{
		"email": insertedDoc.Email,
		"password": dataUser.Password,
		"salt": dataUser.Salt,
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

// func UpdateUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
// 	dataUser, err := GetUserFromID(iduser, db)
// 	if err != nil {
// 		return err
// 	}
// 	if insertedDoc.Email == "" || insertedDoc.Password == "" {
// 		return fmt.Errorf("Dimohon untuk melengkapi data")
// 	}
// 	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
// 		return fmt.Errorf("Email tidak valid")
// 	}
// 	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
// 	if existsDoc.Email == insertedDoc.Email {
// 		return fmt.Errorf("Email sudah terdaftar")
// 	}
// 	if strings.Contains(insertedDoc.Password, " ") {
// 		return fmt.Errorf("password tidak boleh mengandung spasi")
// 	}
// 	if len(insertedDoc.Password) < 8 {
// 		return fmt.Errorf("password terlalu pendek")
// 	}
// 	salt := make([]byte, 16)
// 	_, err = rand.Read(salt)
// 	if err != nil {
// 		return fmt.Errorf("kesalahan server : salt")
// 	}
// 	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
// 	user := bson.M{
// 		"email": insertedDoc.Email,
// 		"password": hex.EncodeToString(hashedPassword),
// 		"salt": hex.EncodeToString(salt),
// 		"role": dataUser.Role,
// 	}
// 	err = UpdateOneDoc(iduser, db, "user", user)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func UpdatePasswordUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Password) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	salt, err := hex.DecodeString(dataUser.Salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != dataUser.Password {
		return fmt.Errorf("password lama salah")
	}
	if insertedDoc.Newpassword == ""  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if strings.Contains(insertedDoc.Newpassword, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Newpassword) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt = make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Newpassword), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email": dataUser.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.User) error {
	dataUser, err := GetUserFromID(iduser, db)
	if err != nil {
		return err
	}
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"email": insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"salt": hex.EncodeToString(salt),
		"role": dataUser.Role,
	}
	err = UpdateOneDoc(iduser, db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser(db *mongo.Database) (user []model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser mongo: %s", err)
	}
	err = cursor.All(context.Background(), &user)
	if err != nil {
		return user, fmt.Errorf("error GetAllUser context: %s", err)
	}
	return user, nil
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// pengguna
func UpdatePengguna(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Pengguna) error {
	pengguna, err := GetPenggunaFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if pengguna.ID != idparam {
		return fmt.Errorf("Anda bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.TanggalLahir == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == ""{
		return fmt.Errorf("Dimohon untuk melengkapi data")
	} 
	pgn := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"tanggallahir": insertedDoc.TanggalLahir,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"akun": model.User {
			ID : pengguna.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "pengguna", pgn)
	if err != nil {
		return err
	}
	return nil
}

func GetAllPengguna(db *mongo.Database) (pengguna []model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func GetPenggunaFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetPenggunaFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("pengguna tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//by admin
func GetPenggunaFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (pengguna model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	user, err := GetUserFromID(pengguna.Akun.ID, db)
	if err != nil {
		return pengguna, fmt.Errorf("error GetPenggunaFromID mongo: %s", err)
	}
	akun := model.User{
		ID: user.ID,
		Email: user.Email,
		Role: user.Role,
	}
	pengguna.Akun = akun
	return pengguna, nil
}

func GetAllPenggunaByAdmin(db *mongo.Database) (pengguna []model.Pengguna, err error) {
	collection := db.Collection("pengguna")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna mongo: %s", err)
	}
	err = cursor.All(context.Background(), &pengguna)
	if err != nil {
		return pengguna, fmt.Errorf("error GetAllPengguna context: %s", err)
	}
	return pengguna, nil
}

func GetDriverFromIDByAdmin(idparam primitive.ObjectID, db *mongo.Database) (driver model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{
		"_id": idparam,
	}
	err = collection.FindOne(context.Background(), filter).Decode(&driver)
	if err != nil {
		return driver, err
	}
	user, err := GetUserFromID(driver.Akun.ID, db)
	if err != nil {
		return driver, err
	}
	akun := model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	driver.Akun = akun
	return driver, nil
}

// func GetAllDriverByAdmin(db *mongo.Database) (driver model.Driver, err error) {
// 	collection := db.Collection("driver")
// 	filter := bson.M{}
// 	cursor, err := collection.Find(context.Background(), filter)
// 	if err != nil {
// 		return driver, err
// 	}
// 	err = cursor.All(context.Background(), &driver)
// 	if err != nil {
// 		return driver, err
// 	}
// 	for _, m := range driver {
// 		user, err := GetUserFromID(m.Akun.ID, db)
// 		if err != nil {
// 			return driver, err
// 		}
// 		akun := model.User{
// 			ID:    user.ID,
// 			Email: user.Email,
// 			Role:  user.Role,
// 		}
// 		m.Akun = akun
// 		driver = append(driver, m)
// 		driver = driver[1:]
// 	}
// 	return driver, nil
// }


// driver
func UpdateDriver(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Driver) error {
	driver, err := GetDriverFromAkun(iduser, db)
	if err != nil {
		return err
	}
	if driver.ID != idparam {
		return fmt.Errorf("anda bukan pemilik data ini")
	}
	if insertedDoc.NamaLengkap == "" || insertedDoc.JenisKelamin == "" || insertedDoc.NomorHP == "" || insertedDoc.Alamat == "" || insertedDoc.PlatBis == ""{
		return fmt.Errorf("dimohon untuk melengkapi data")
	} 
	drv := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"jeniskelamin": insertedDoc.JenisKelamin,
		"nomorhp": insertedDoc.NomorHP,
		"alamat": insertedDoc.Alamat,
		"platbis": insertedDoc.PlatBis,
		"akun": model.User {
			ID : driver.Akun.ID,
		},
	}
	err = UpdateOneDoc(idparam, db, "driver", drv)
	if err != nil {
		return err
	}
	return nil
}

func GetAllDriver(db *mongo.Database) (driver []model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return driver, fmt.Errorf("error GetAllDriver mongo: %s", err)
	}
	err = cursor.All(context.Background(), &driver)
	if err != nil {
		return driver, fmt.Errorf("error GetAllDriver context: %s", err)
	}
	return driver, nil
}

func GetDriverFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetDriverFromAkun(akun primitive.ObjectID, db *mongo.Database) (doc model.Driver, err error) {
	collection := db.Collection("driver")
	filter := bson.M{"akun._id": akun}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("Driver tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//obat

func InsertTiket(iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Tiket) error {
	if insertedDoc.TujuanEvent == "" || insertedDoc.Jemputan == "" || insertedDoc.Keterangan == "" || insertedDoc.Harga == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}

	tkt := bson.M{
		"tujuaneven":    insertedDoc.TujuanEvent,
		"jemputan":   insertedDoc.Jemputan,
		"keterangan":   insertedDoc.Keterangan,
		"harga":        insertedDoc.Harga,
	}

	_, err := InsertOneDoc(db, "tiket", tkt)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data tiket: %s", err)
	}
	return nil
}

func UpdateTiket(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.Tiket) error {
	_, err := GetTiketFromID(idparam, db)
	if err != nil {
		return err
	}
	if insertedDoc.TujuanEvent == "" || insertedDoc.Jemputan == "" || insertedDoc.Keterangan == "" || insertedDoc.Harga == ""  {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	tkt := bson.M{
		"tujuaneven":    insertedDoc.TujuanEvent,
		"jemputan":   insertedDoc.Jemputan,
		"keterangan":   insertedDoc.Keterangan,
		"harga":        insertedDoc.Harga,
	}

	err = UpdateOneDoc(idparam, db, "tiket", tkt)
	if err != nil {
		return err
	}
	return nil
}


func DeleteTiket(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetTiketFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "tiket")
	if err != nil {
		return err
	}
	return nil
}

func GetAllTiket(db *mongo.Database) (tiket []model.Tiket, err error) {
	collection := db.Collection("tiket")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return tiket, fmt.Errorf("error GetAllTiket mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &tiket)
	if err != nil {
		return tiket, fmt.Errorf("error GetAllTiket context: %s", err)
	}
	return tiket, nil
}


func GetTiketFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Tiket, err error) {
	collection := db.Collection("tiket")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

//order

func InsertOrderTiket(idparam, iduser primitive.ObjectID, db *mongo.Database, insertedDoc model.OrderTiket) error {

	if insertedDoc.Event == "" || insertedDoc.Quantity == "" || insertedDoc.TotalCost == "" || insertedDoc.Status == "" {
		return fmt.Errorf("harap lengkapi semua data order")
	}

	ord := bson.M{
		"pengguna": bson.M{
			"_id" : iduser,
		},
		"driver": bson.M{
			"_id" : insertedDoc.Driver.ID,
		},
		"tiket": bson.M{
			"_id" : idparam,
		},
		"even":    insertedDoc.Event,
		"quantity":    insertedDoc.Quantity,
		"total_cost":   insertedDoc.TotalCost,
		"status":   insertedDoc.Status,
	}

	_, err := InsertOneDoc(db, "order", ord)
	if err != nil {
		return fmt.Errorf("error saat menyimpan data order tiket: %s", err)
	}
	return nil
}


//update status pengiriman
func UpdateStatusOrderTiket(idorder primitive.ObjectID, db *mongo.Database, insertedDoc model.OrderTiket) error {
	order, err := GetOrderFromID(idorder, db)
	if err != nil {
		return err
	}

	data := bson.M{
		"even":    order.Event,
		"quantity":    order.Quantity,
		"total_cost":   order.TotalCost,
		"status": insertedDoc.Status,
	}

	err = UpdateOneDoc(idorder, db, "order", data)
	if err != nil {
		return err
	}
	return nil
}


func DeleteOrder(idparam, iduser primitive.ObjectID, db *mongo.Database) error {
	_, err := GetOrderFromID(idparam, db)
	if err != nil {
		return err
	}
	err = DeleteOneDoc(idparam, db, "order")
	if err != nil {
		return err
	}
	return nil
}

func GetOrderFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.OrderTiket, err error) {
	collection := db.Collection("order")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, err
	}
	return doc, nil
}


func GetAllOrder(db *mongo.Database) (order []model.OrderTiket, err error) {
	collection := db.Collection("order")
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder mongo: %s", err)
	}
	err = cursor.All(context.TODO(), &order)
	if err != nil {
		return order, fmt.Errorf("error GetAllOrder context: %s", err)
	}
	return order, nil
}
