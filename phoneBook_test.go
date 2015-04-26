package main

import (
	"github.com/julienschmidt/httprouter"
	//"net/http"
	//"net/http/httptest"
	//"github.com/Charlesworth/phoneBook"
	//"fmt"
	"github.com/boltdb/bolt"
	"reflect"
	"runtime"
	"testing"
)

//func main() {}

func TestNewRouter(t *testing.T) {
	//Get NewRouter to return *httpRouter
	router := NewRouter()

	//Did it return anything?
	if router == nil {
		t.Error("Router() did not return anything *httprouter.Router")
	}

	//did it return a *httpRouter?
	if reflect.TypeOf(router) != reflect.TypeOf(httprouter.New()) {
		t.Error("Router() did not return type *httprouter")
	}

}

func TestSetProc(t *testing.T) {
	SetProc()

	if runtime.GOMAXPROCS(0) != 2 {
		t.Error("Application not using 2 processors as set by setProc()")
	}
}

func Testmain(t *testing.T) {

}

func TestNewBoltClient(t *testing.T) {
	boltTest := NewBoltClient()

	boltTest.Mutex.Lock()

	//Test creating bucket
	err := boltTest.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("test"))

		return err
	})
	if err != nil {
		t.Error(err)
	}

	//Write to bucket
	err = boltTest.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		err = b.Put([]byte("hello"), []byte("world"))
		return err
	})
	if err != nil {
		t.Error(err)
	}

	boltTest.Mutex.Unlock()
	boltTest.Mutex.RLock()

	//Get from bucket
	var v []byte
	boltTest.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		v = b.Get([]byte("hello"))
		return nil
	})
	if v == nil {
		t.Error("Cannot retrive values from BOLTDB")
	}
	if string(v) != "world" {
		t.Error("BoltDB not storing values correctly")
	}

	boltTest.Mutex.RUnlock()
	boltTest.Mutex.Lock()

	//Delete from bucket
	err = boltTest.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		err = b.Delete([]byte("hello"))
		return err
	})
	if err != nil {
		t.Error(err)
	}

	//Delete test bucket
	err = boltTest.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("test"))
		return err
	})
	if err != nil {
		t.Error(err)
	}

	boltTest.Mutex.Unlock()

	boltTest.DB.Close()
}
