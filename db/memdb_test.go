package db_test

import (
	"testing"

	"github.com/getircase/db"
	"github.com/stretchr/testify/require"
)

func TestInMemDbError(t *testing.T) {
	// no value provided in the value

	err := db.MemDBMgr().SetKV("test", "")
	require.NoError(t, err)
	// Unknown key
	// err no key found
	_, err = db.MemDBMgr().Retrieve("test1")
	require.NotEmpty(t, err)
	require.Errorf(t, err, "Key not found")

	// Set the value with key will be errored
	err = db.MemDBMgr().SetKV("", "testValue")
	require.NotEmpty(t, err)
	require.Errorf(t, err, "Key cannot be empty")

}

func TestInMemDb(t *testing.T) {
	k, v := "test", "testValue"
	var out interface{}
	// initialize the db with key value
	err := db.MemDBMgr().SetKV(k, v)
	require.NoError(t, err)
	// err no key found
	out, err = db.MemDBMgr().Retrieve(k)
	require.Empty(t, err)
	rs, ok := out.(map[string]string)
	require.Equal(t, true, ok)
	require.NotNil(t, rs)
	require.Equal(t, rs["key"], k)
	require.Equal(t, rs["value"], v)
}
