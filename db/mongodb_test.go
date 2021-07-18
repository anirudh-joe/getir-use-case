package db_test

import (
	"testing"

	"github.com/getircase/db"
	"github.com/getircase/models"
	"github.com/stretchr/testify/require"
)

func TestMongoDbError(t *testing.T) {
	var req models.MongoRequest
	var resp models.MongoResponse
	req.EndDate = "2016-03-02"
	req.MinCount = 2900
	req.MaxCount = 3000
	// start date parse error Month out of range
	req.StartDate = "2016-13-26"
	_, err := db.MongoMgr().Retrieve(req)
	require.NotEmpty(t, err)
	require.Errorf(t, err, "parsing time \"2016-13-26\": month out of range")
	req.StartDate = "2016-01-32"
	// start date parse error day out of range
	_, err = db.MongoMgr().Retrieve(req)
	require.NotEmpty(t, err)
	require.Errorf(t, err, "parsing time \"2016-01-32\": day out of range")

	// end date parse error Month out of range
	req.EndDate = "2018-13-26"
	req.StartDate = "2016-01-26"
	_, err = db.MongoMgr().Retrieve(req)
	require.NotEmpty(t, err)
	require.Errorf(t, err, "parsing time \"2018-13-26\": month out of range")

	// end date parse error day out of range
	req.EndDate = "2018-01-32"
	req.StartDate = "2016-01-26"
	_, err = db.MongoMgr().Retrieve(req)
	require.NotEmpty(t, err)
	require.Errorf(t, err, "parsing time \"2018-13-26\": month out of range")

	req.StartDate = "2016-01-02"
	req.EndDate = "2016-03-02"
	req.MinCount = 3100
	req.MaxCount = 3000
	rs, err := db.MongoMgr().Retrieve(req)
	require.NotEmpty(t, err)
	require.Errorf(t, err, "no data found")
	require.NotNil(t, rs)

	resp = rs.(models.MongoResponse)
	require.Equal(t, 204, resp.Code)
	require.Equal(t, "No Data Found", resp.Msg)
	require.Equal(t, 0, len(resp.Records))

}

func TestMongoDb(t *testing.T) {
	var req models.MongoRequest
	var resp models.MongoResponse
	req.StartDate = "2016-01-02"
	req.EndDate = "2016-03-02"
	req.MinCount = 2900
	req.MaxCount = 3000
	rs, err := db.MongoMgr().Retrieve(req)
	require.Empty(t, err)
	require.NotNil(t, rs)

	resp = rs.(models.MongoResponse)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "Success", resp.Msg)
	require.LessOrEqual(t, 0, len(resp.Records))

}
