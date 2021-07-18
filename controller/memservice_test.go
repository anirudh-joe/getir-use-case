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
)

func TestMemDbHandlerGetError(t *testing.T) {
	var memServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	// key request Parameter not Provided Error Case

	// Create a request to pass to our handler. We don't have any query parameters for GET, so
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("GET", "/in-memory", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	memServer = new(controller.MemDbGate)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusForbidden)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("key request Parameter not Provided"))

	// Key not found Error Case
	// Create another request as similar to before with no request body but add the path parameter for 'key'
	req, err = http.NewRequest("GET", "/in-memory?key=test", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusNotFound)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("Key not found"))
}

func TestMemDbHandlerPostError(t *testing.T) {

	var memServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	rq := map[string]string{}
	// Body wiith Empty key
	memServer = new(controller.MemDbGate)
	// Generate Body object to pass to request
	rq["key"] = ""
	rq["value"] = "testValue"
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/in-memory", reader)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("Key cannot be empty"))

	// Nil Body Case
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/in-memory", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("No request content to process"))
}

func TestMemDbHandlerPostModelError(t *testing.T) {
	var memServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	rq := map[interface{}]interface{}{}
	// Bad Request Body Case
	rq["abc"] = 1.1
	memServer = new(controller.MemDbGate)
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/in-memory", reader)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusInternalServerError)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.Equal(t, body, []byte("unexpected end of JSON input"))
}

func TestMemDbHandlerPostSuccess(t *testing.T) {
	var memServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var resp models.InMemory
	rq := map[string]string{}
	// Store data in in-memdb
	memServer = new(controller.MemDbGate)
	rq["key"] = "test"
	rq["value"] = "testValue"
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	// Create a Http.POST request to pass to our handler.
	req, err = http.NewRequest("POST", "/in-memory", reader)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusCreated)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &resp)
	// Check the Response Body and it is what we expect.
	require.NoError(t, err)
	require.Equal(t, resp.Key, rq["key"])
	require.Equal(t, resp.Value, rq["value"])
}

func TestMemDbHandlerGetSuccess(t *testing.T) {
	var memServer http.Handler
	var req *http.Request
	var err error
	var body []byte
	var resp models.InMemory
	rq := map[string]string{}
	// Store the data first and fetch
	memServer = new(controller.MemDbGate)
	rq["key"] = "test"
	rq["value"] = "testValue"
	payloadBytes, _ := json.Marshal(rq)
	reader := bytes.NewReader(payloadBytes)
	req, err = http.NewRequest("POST", "/in-memory", reader)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	memServer.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusCreated)
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	require.Equal(t, resp.Key, rq["key"])
	require.Equal(t, resp.Value, rq["value"])
	// Successful get for the above stored key
	// Create another request with no request body but add the path parameter for 'key' for the stored key
	req, err = http.NewRequest("GET", "/in-memory?key=test", nil)
	require.NoError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	memServer = new(controller.MemDbGate)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	memServer.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	require.Equal(t, rr.Code, http.StatusAccepted)
	// Read the response body
	body, err = ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	require.Equal(t, resp.Key, rq["key"])
	require.Equal(t, resp.Value, rq["value"])

}
