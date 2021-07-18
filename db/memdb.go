package db

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

// MemDBManager ...
// Interface Pattern with following functions SetKV Retrieve
type MemDBManager interface {
	SetKV(key, value string) error
	Retrieve(key string) (out interface{}, err error)
}

// memdb ...
// Unexported memdb object for not be misused
// SingleTon Pattern
type memdb struct {
	db *badger.DB
}

// inMemInstance ...
// By default, Badger ensures all the data is persisted to the disk.When Badger is running in in-memory mode
// All the data is stored in the memory. Reads and writes are much faster in in-memory mode,
// but all the data stored in Badger will be lost in case of a crash or close.
// To open badger in in-memory mode, set the InMemory option.
func inMemInstance() MemDBManager {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	return &memdb{db: db}

}

// Unexported object of the Interface
var memdbmgr = inMemInstance()

// MemDBMgr ...
// Exported function to be consumed from anywhere in the project
// Return the Interface instance to expose the underlying functionality
func MemDBMgr() MemDBManager { return memdbmgr }

// SetKV ...
// Set key and associated value for the badger`s in-memory db
func (m *memdb) SetKV(key, value string) error {
	err := m.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
		err := txn.SetEntry(e)
		return err
	})
	return err
}

// Retrieve ...
// fetches the associated data from In-memory badger DB for the requested key
// return err if Key not found or Key is empty
func (m *memdb) Retrieve(key string) (out interface{}, err error) {
	var valCopy []byte
	err = m.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(key))
		if e != nil {
			return e
		}
		valCopy, e = item.ValueCopy(nil)
		if e != nil {
			return e
		}
		return nil
	})
	if err == badger.ErrKeyNotFound || err == badger.ErrEmptyKey {
		return
	}
	rs := map[string]string{}
	rs["key"] = key
	rs["value"] = string(valCopy)
	out = rs
	return
}
