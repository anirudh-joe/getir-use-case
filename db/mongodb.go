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

type MongoManager interface {
	Retrieve(input interface{}) (out interface{}, err error)
}

type mongodb struct {
	collection *mongo.Collection
}

// mongoInstance ...
// By default, Badger ensures all the data is persisted to the disk.When Badger is running in in-memory mode
// All the data is stored in the memory. Reads and writes are much faster in in-memory mode,
// but all the data stored in Badger will be lost in case of a crash or close.
// To open badger in in-memory mode, set the InMemory option.
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
	// get cursor for generate pipeline
	cursor, err := m.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
	// defer the closing of the cursor
	defer cursor.Close(context.TODO())
	// read all data from cursor result to records Data
	if err = cursor.All(context.TODO(), &recordsData); err != nil {
		mr.Msg = err.Error()
		return mr, err
	}
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
