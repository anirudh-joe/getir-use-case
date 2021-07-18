package models

// MongoRequest ...
// Model for MongoRequest http requests
type MongoRequest struct {
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	MinCount  float64 `json:"minCount"`
	MaxCount  float64 `json:"maxCount"`
}
