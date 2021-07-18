package models

import "go.mongodb.org/mongo-driver/bson"

// MongoResponse ...
// Model for MongoDb http response
type MongoResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []bson.M `json:"records"`
}
