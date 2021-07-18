package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/getircase/db"
	"github.com/getircase/models"
)

// MemDbGate ...
// In-Memory Gateway/handler for in-memory based request
type MemDbGate struct{}

// ServeHTTP ...
// Generic ServeHttp linked with MemDbGate
// Serves HTTP method GET and POST
// uri path value /in-memory
func (gate *MemDbGate) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	var err error
	var result interface{}
	var out, body []byte

	if request.Method == "GET" {
		// Get Key path parameter from url path
		keys, ok := request.URL.Query()["key"]
		// No found the required path throw http.StatusForbidden
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte("key request Parameter not Provided"))
			return
		}
		// Not found or more than required number of parameters in the url path throw http.StatusBadRequest
		if len(keys) != 1 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("key request Parameter not Provided"))
			return
		}
		// Retrieve associated Value for the requested key
		// If there is an error throw http.StatusNotFound
		result, err = db.MemDBMgr().Retrieve(keys[0])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		// return Associated Value for the requested key in the requested format
		gr := result.(map[string]string)
		out, _ = json.Marshal(gr)
		rw.WriteHeader(http.StatusAccepted)
		rw.Write(out)
	} else if request.Method == "POST" {
		// If body not present throw http.StatusInternalServerError
		if nil == request.Body {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte("No request content to process"))
			return
		}
		defer request.Body.Close()
		// If body can not be read throw http.StatusInternalServerError
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		var content models.InMemory
		// Map request body to not InMemory model
		// otherwise throw http.StatusInternalServerError
		if err = json.Unmarshal(body, &content); err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		// Set the associated value for the given key
		// if err is present throw http.StatusInternalServerError
		err = db.MemDBMgr().SetKV(content.Key, content.Value)
		if err != nil {
			rw.WriteHeader(500)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		// otherwise return code http.StatusCreated
		// write response for the request
		rw.WriteHeader(http.StatusCreated)
		rw.Write(body)
	}
}
