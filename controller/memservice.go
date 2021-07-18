package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/getircase/db"
	"github.com/getircase/models"
)

type MemDbGate struct {
	// handler http.Handler
}

// ServeHTTP ...
func (gate *MemDbGate) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	var err error
	var result interface{}
	var out, body []byte

	if request.Method == "GET" {
		keys, ok := request.URL.Query()["key"]
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte("key request Parameter not Provided"))
			return
		}
		if len(keys) != 1 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("key request Parameter not Provided"))
			return
		}
		result, err = db.MemDBMgr().Retrieve(keys[0])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		gr := result.(map[string]string)
		out, _ = json.Marshal(gr)
		rw.WriteHeader(http.StatusAccepted)
		rw.Write(out)
	} else if request.Method == "POST" {
		if nil == request.Body {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte("No request content to process"))
			return
		}
		defer request.Body.Close()
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		var content models.InMemory
		if err = json.Unmarshal(body, &content); err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}

		err = db.MemDBMgr().SetKV(content.Key, content.Value)
		if err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusCreated)
		rw.Write(body)
	}
}
