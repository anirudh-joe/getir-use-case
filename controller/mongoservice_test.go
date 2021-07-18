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
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	mongoServer = new(controller.MongoDbGate)
	// Bad Http Method Case
	req, err = http.NewRequest("PUT", "/mongo", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "http Method not supported")
	require.Equal(t, mr.Records, []bson.M{})

}
func TestMongoHandlerNilBodyError(t *testing.T) {
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	mongoServer = new(controller.MongoDbGate)
	// Nil body Error Case
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/mongo", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "No request content to process")
	require.Equal(t, mr.Records, []bson.M{})
}

func TestMongoHandlerModelError(t *testing.T) {
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	rq := map[interface{}]interface{}{}
	// Bad Request Body Case
	rq["abc"] = 1.1
	mongoServer = new(controller.MongoDbGate)
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, 400)
	require.Equal(t, mr.Msg, "unexpected end of JSON input")
	require.Equal(t, mr.Records, []bson.M{})
}

func TestMongoHandlerResponseError(t *testing.T) {
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var mr models.MongoResponse
	rq := map[string]interface{}{}
	// Bad Request Object case
	// Generate Body object to pass to request
	rq["startDate"] = "2016-01-32"
	rq["endDate"] = "2016-03-02"
	rq["minCount"] = 3100
	rq["maxCount"] = 3000
	mongoServer = new(controller.MongoDbGate)
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusNotFound)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, http.StatusBadRequest)
	require.Equal(t, mr.Msg, "parsing time \"2016-01-32\": day out of range")
	require.Equal(t, mr.Records, []bson.M{})

	// No records found for the requested filters
	rq["startDate"] = "2016-01-02"
	payloadBytes, _ = json.Marshal(rq)
	reader = bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusNotFound)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &mr)
	require.NoError(t, err)
	require.Equal(t, mr.Code, http.StatusNoContent)
	require.Equal(t, mr.Msg, "No Data Found")
	require.Equal(t, mr.Records, []bson.M{})
}

func TestMongoDbHandlerPostSuccess(t *testing.T) {
	var mongoServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var resp models.MongoResponse
	rq := map[string]interface{}{}
	// Generate Body object to pass to request
	mongoServer = new(controller.MongoDbGate)
	rq["startDate"] = "2016-01-02"
	rq["endDate"] = "2016-06-02"
	rq["minCount"] = 2900
	rq["maxCount"] = 3000
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/mongo", reader)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mongoServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusAccepted)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	require.Equal(t, resp.Code, 0)
	require.Equal(t, resp.Msg, "Success")
	require.LessOrEqual(t, 0, len(resp.Records))
}
