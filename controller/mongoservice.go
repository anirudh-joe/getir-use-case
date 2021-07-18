package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/getircase/db"
	"github.com/getircase/models"
	"go.mongodb.org/mongo-driver/bson"
)

// MongoDbGate ...
// In-Memory Gateway/handler for in-memory based request
type MongoDbGate struct{}

// ServeHTTP ...
// Generic ServeHttp linked with MongodbGate
// Serves HTTP method POST
// uri path /mongo
func (gate *MongoDbGate) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	var result interface{}
	var out []byte
	var mr models.MongoResponse
	data := []bson.M{}

	mr.Code = http.StatusBadRequest
	mr.Records = data
	// If request method is not POST throw http.StatusInternalServerError
	if request.Method != "POST" {
		mr.Msg = "http Method not supported"
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	// If body not present throw http.StatusInternalServerError
	if nil == request.Body {
		mr.Msg = "No request content to process"
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}

	defer request.Body.Close()
	// If body can not be read throw http.StatusInternalServerError
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		mr.Msg = err.Error()
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	var content models.MongoRequest
	// Map request body to not MongoRequest model
	// otherwise throw http.StatusInternalServerError
	if err = json.Unmarshal(body, &content); err != nil {
		rw.WriteHeader(500)
		mr.Msg = err.Error()
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	// Retrieve associated data for the requested filters
	// If there is an error throw http.StatusNotFound
	result, err = db.MongoMgr().Retrieve(content)
	out, _ = json.Marshal(result)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		_, _ = rw.Write(out)
		return
	}
	// otherwise return code http.StatusAccepted
	// write response for the request
	rw.WriteHeader(http.StatusAccepted)
	_, _ = rw.Write(out)
}
