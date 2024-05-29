package minidb

import (
	"errors"
	"sync"
)

type MiniDB struct {
	mu   sync.RWMutex
	data map[string]string
	rev  map[string]string
}

func New() *MiniDB {
	return &MiniDB{
		mu:   sync.RWMutex{},
		data: make(map[string]string),
		rev:  make(map[string]string),
	}
}

func (db *MiniDB) Get(key string) (val string, err error) {
	db.mu.RLock()
	if v, ok := db.data[key]; ok {
		val = v
	} else {
		err = errors.New("key not found")
	}
	db.mu.RUnlock()

	return val, err
}

func (db *MiniDB) Resolve(val string) (key string, err error) {
	db.mu.RLock()
	if k, ok := db.rev[val]; ok {
		key = k
	} else {
		err = errors.New("value not found")
	}
	db.mu.RUnlock()

	return key, err
}

func (db *MiniDB) Set(key string, val string) (err error) {
	db.mu.Lock()
	if _, ok := db.data[key]; !ok {
		if _, ok := db.rev[val]; !ok {
			db.data[key] = val
			db.rev[val] = key
		} else {
			err = errors.New("value already exist")
		}
	} else {
		err = errors.New("key already exist")
	}
	db.mu.Unlock()

	return err
}
