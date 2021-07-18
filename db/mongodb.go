package db

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getircase/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoManager ...
// Interface Pattern with following functions Retrieve
type MongoManager interface {
	Retrieve(input interface{}) (out interface{}, err error)
}

// mongodb ...
// Unexported memdb object for not be misused
// SingleTon Pattern
type mongodb struct {
	collection *mongo.Collection
}

// mongoInstance ...
// driver starts with creating a Client from a connection string and subsequently get session
// Database and Collection types can be used to access the database and collection
// return the collection embedded inside the unexported mongodb struct
func mongoInstance() MongoManager {

	context, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(context, options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"))
	if err != nil {
		panic(err)
	}
	ss, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	database := ss.Client().Database("getir-case-study")

	recordsCollection := database.Collection("records")
	return &mongodb{collection: recordsCollection}
}

var mongomgr = mongoInstance()

func MongoMgr() MongoManager { return mongomgr }

// Retrieve ...
// Implemets mongodb aggregate functionality with Pipeline stages
// first stage is to match the records createdAt between start and end date
// second stage is to project the records and take the sum of count array and put inside totalCount
// third stage is to match the records totalCount between minCount and maxCount
func (m *mongodb) Retrieve(input interface{}) (out interface{}, err error) {
	recordsData := []bson.M{}
	var mr models.MongoResponse
	var req models.MongoRequest

	req, _ = input.(models.MongoRequest)

	mr.Code = http.StatusBadRequest
	mr.Records = recordsData
	// Convert incoming date in request to epoch time
	sd, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
	// Convert incoming date in request to epoch time
	ed, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
	// pipeline parameter must be an array of documents, each representing an aggregation stage.
	//  Documents pass through the stages in sequence.
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gt": sd,
					"$lt": ed,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"key":        1,
				"createdAt":  1,
				"totalCount": bson.M{"$sum": "$counts"},
			},
		},
		{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gt": req.MinCount,
					"$lt": req.MaxCount,
				},
			},
		},
	}
	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	cursor, err := m.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
	// defer the closing of the cursor
	defer cursor.Close(context.TODO())
	// Cursor.All will decode all of the returned elements at once
	if err = cursor.All(context.TODO(), &recordsData); err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
	// return the data if the records are found after pipeline executes
	if len(recordsData) > 0 {
		mr.Code = 0
		mr.Msg = "Success"
		mr.Records = recordsData
		return mr, nil
	}

	mr.Code = http.StatusNoContent
	mr.Msg = "No Data Found"
	mr.Records = recordsData
	err = fmt.Errorf("no data found")
	return mr, err
}
