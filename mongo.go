package gobd

import (
	"context"
	"os"
	"time"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetConnectionMongo(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

func GetAllData(MongoConnect *mongo.Database, colname string) []bisdata {
	data := atdb.GetAllDoc[[]bisdata](MongoConnect, colname)
	return data
}

// func InsertDataLonlat(MongoConn *mongo.Database, colname string, coordinate []float64, name, volume, tipe string) (InsertedID interface{}) {
// 	req := new(LonLatProperties)
// 	req.Type = tipe
// 	req.Coordinates = coordinate
// 	req.Name = name
// 	req.Volume = volume

// 	ins := atdb.InsertOneDoc(MongoConn, colname, req)
// 	return ins
// }

func InsertDatabisini(MongoConn *mongo.Database, colname string, jamgo	time.Time,	jamout time.Time, jemputan string, nokursi string) (InsertedID interface{}) {
	req := new(Bis)
	req.NoKursi = nokursi
	req.Jemputan = jemputan
	req.JamGo =	jamgo 
	req.JamOut = jamout

	ins := atdb.InsertOneDoc(MongoConn, colname, req)
	return ins
}

func UpdateDataGeojson(MongoConn *mongo.Database, colname, name, newVolume, newTipe string) error {
    // Filter berdasarkan nama
    filter := bson.M{"name": name}

    // Update data yang akan diubah
    update := bson.M{
        "$set": bson.M{
            "volume": newVolume,
            "tipe":   newTipe,
        },
    }

    // Mencoba untuk mengupdate dokumen
    _, err := MongoConn.Collection(colname).UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

func DeleteDataGeojson(MongoConn *mongo.Database, colname string, name string) (*mongo.DeleteResult, error) {
    filter := bson.M{"name": name}
    del, err := MongoConn.Collection(colname).DeleteOne(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    return del, nil
}