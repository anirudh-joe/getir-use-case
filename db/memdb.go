package db

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

type MemDBManager interface {
	SetKV(key, value string) error
	Retrieve(key string) (out interface{}, err error)
}

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

var memdbmgr = inMemInstance()

func MemDBMgr() MemDBManager { return memdbmgr }

func (m *memdb) SetKV(key, value string) error {
	err := m.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
		err := txn.SetEntry(e)
		return err
	})
	return err
}

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
