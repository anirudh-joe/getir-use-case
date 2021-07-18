package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/getircase/db"
	"github.com/getircase/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDbGate struct {
	// handler http.Handler
}

func (gate *MongoDbGate) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	var result interface{}
	var out []byte
	var mr models.MongoResponse
	data := []bson.M{}

	mr.Code = http.StatusBadRequest
	mr.Records = data
	if request.Method != "POST" {
		mr.Msg = "http Method not supported"
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	if nil == request.Body {
		mr.Msg = "No request content to process"
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		mr.Msg = err.Error()
		rw.WriteHeader(500)
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}
	var content models.MongoRequest
	if err = json.Unmarshal(body, &content); err != nil {
		rw.WriteHeader(500)
		mr.Msg = err.Error()
		out, _ = json.Marshal(mr)
		_, _ = rw.Write(out)
		return
	}

	result, err = db.MongoMgr().Retrieve(content)
	out, _ = json.Marshal(result)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		_, _ = rw.Write(out)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	_, _ = rw.Write(out)
}
