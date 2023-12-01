package gobd

import "time"

// type GeometryPolygon struct {
// 	Coordinates [][][]float64 `json:"coordinates" bson:"coordinates"`
// 	Type        string        `json:"type" bson:"type"`
// }

// type GeometryLineString struct {
// 	Coordinates [][]float64 `json:"coordinates" bson:"coordinates"`
// 	Type        string      `json:"type" bson:"type"`
// }

// type GeometryPoint struct {
// 	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
// 	Type        string    `json:"type" bson:"type"`
// }

// type GeoJsonLineString struct {
// 	Type       string             `json:"type" bson:"type"`
// 	Properties Properties         `json:"properties" bson:"properties"`
// 	Geometry   GeometryLineString `json:"geometry" bson:"geometry"`
// }

// type GeoJsonPolygon struct {
// 	Type       string          `json:"type" bson:"type"`
// 	Properties Properties      `json:"properties" bson:"properties"`
// 	Geometry   GeometryPolygon `json:"geometry" bson:"geometry"`
// }

// type Geometry struct {
// 	Coordinates interface{} `json:"coordinates" bson:"coordinates"`
// 	Type        string      `json:"type" bson:"type"`
// }
// type GeoJson struct {
// 	Type       string     `json:"type" bson:"type"`
// 	Properties Properties `json:"properties" bson:"properties"`
// }

 type Properties struct {
 	NamaBis string `json:"name" bson:"name"`
 }

// type LonLatProperties struct {
// 	Type        string    `json:"type" bson:"type"`
// 	Name        string    `json:"name" bson:"name"`
// 	Volume      string    `json:"volume" bson:"volume"`
// 	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
// }

 type bisdata struct {
 	Type       string     `json:"type" bson:"type"`
	Properties Properties `json:"properties" bson:"properties"`
 }

type Credents struct {
	Status  string `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message" bson:"message"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	// Email		 string             	`bson:"email,omitempty" json:"email,omitempty"`
}
type Admin struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`

}

type Shuttle struct {
	Username string `json:"username" bson:"username"`
	jamgo 	string `json:"password" bson:"password"`
	jamout     string `json:"role" bson:"role"`
	Role     string `json:"role" bson:"role"`
}

type Bis struct {
	NoKursi string `json:"nokursi" bson:"nokursi"`
	Jemputan string `json:"jemputan" bson:"jemputan"`
	JamGo	time.Time	`json:"jamgo" bson:"jamgo"`
    JamOut time.Time `json:"jamout" bson:"jamout"`
}