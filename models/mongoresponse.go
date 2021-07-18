package models

import "go.mongodb.org/mongo-driver/bson"

// MongoResponse ...
type MongoResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []bson.M `json:"records"`
}
