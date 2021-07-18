package models

// InMemory ...
// Model for MongoDb http requests
type InMemory struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
