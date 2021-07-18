package controller_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getircase/controller"
	"github.com/getircase/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoHandlerHttpMethodError(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	mongoServer = new(controller.MongoDbGate)

	req, err = http.NewRequest("PUT", "/mongo", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "http Method not supported")
	require.Equal(t, mr.Records, []bson.M{})

}
func TestMongoHandlerNilBodyError(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	mongoServer = new(controller.MongoDbGate)
	req, err = http.NewRequest("POST", "/mongo", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "No request content to process")
	require.Equal(t, mr.Records, []bson.M{})
}

func TestMongoHandlerModelError(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse

	rq := map[interface{}]interface{}{}
	rq["abc"] = 1.1
	mongoServer = new(controller.MongoDbGate)
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)

	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusInternalServerError)

	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &mr)

	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "unexpected end of JSON input")
	require.Equal(t, mr.Records, []bson.M{})

}

func TestMongoHandlerResponseError(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse

	rq := map[string]interface{}{}

	rq["startDate"] = "2016-01-32"
	rq["endDate"] = "2016-03-02"
	rq["minCount"] = 3100
	rq["maxCount"] = 3000

	mongoServer = new(controller.MongoDbGate)
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)

	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusNotFound)

	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &mr)

	require.NoError(t, err)
	require.Equal(t, mr.Code, http.StatusBadRequest)
	require.Equal(t, mr.Msg, "parsing time \"2016-01-32\": day out of range")
	require.Equal(t, mr.Records, []bson.M{})

	rq["startDate"] = "2016-01-02"
	payloadBytes, _ = json.Marshal(rq)
	reader = bytes.NewReader(payloadBytes)

	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusNotFound)

	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &mr)

	require.NoError(t, err)
	require.Equal(t, mr.Code, http.StatusNoContent)
	require.Equal(t, mr.Msg, "No Data Found")
	require.Equal(t, mr.Records, []bson.M{})

}

func TestMongoDbHandlerPostSuccess(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var resp models.MongoResponse
	rq := map[string]interface{}{}
	mongoServer = new(controller.MongoDbGate)
	rq["startDate"] = "2016-01-02"
	rq["endDate"] = "2016-06-02"
	rq["minCount"] = 2900
	rq["maxCount"] = 3000

	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)

	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	mongoServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusAccepted)
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	
	require.Equal(t, resp.Code, 0)
	require.Equal(t, resp.Msg, "Success")
	require.LessOrEqual(t, 0, len(resp.Records))

}
