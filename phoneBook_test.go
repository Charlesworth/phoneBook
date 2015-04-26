// Written by Charles Cochrane, 2015
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

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

	if runtime.GOMAXPROCS(0) != 1 {
		t.Error("Application not using 1 processor as set by setProc()")
	}
}

func TestNewBoltClient(t *testing.T) {
	boltTest := NewBoltClient("BoltTest")

	boltTest.Mutex.Lock()

	//Write to bucket
	err := boltTest.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("BoltTest"))
		err := b.Put([]byte("hello"), []byte("world"))
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
		b := tx.Bucket([]byte("BoltTest"))
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
		b := tx.Bucket([]byte("BoltTest"))
		err = b.Delete([]byte("hello"))
		return err
	})
	if err != nil {
		t.Error(err)
	}

	//Delete test bucket
	err = boltTest.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("BoltTest"))
		return err
	})
	if err != nil {
		t.Error(err)
	}

	boltTest.Mutex.Unlock()

	boltTest.DB.Close()
}

type MockHTTP struct {
	urlStr        string
	expectedWCode int
	JsonStr       []byte
}

func TestPutHandler(t *testing.T) {

	//Setup global variables
	BoltClient = NewBoltClient("test")
	Bucket = "test"

	router := httprouter.New()
	router.PUT("/", putEntryHandler)

	inputHTTP := [5]MockHTTP{
		//test case 0: improper formatted JSON
		{"/", 200, []byte(`asfasdfasdf3421432@:L(*)(*&^"!`)},
		//test case 1: more than 1 element being input in "Entries"
		{"/", 200, []byte(`{"Surname":"Smith","Entries":[{"Firstname":"John","TelNo":"1234567890","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"},{"Firstname":"Jane","TelNo":"1234567890","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}`)},
		//test case 2: new Surname
		{"/", 200, []byte(`{"Surname":"Smith","Entries":[{"Firstname":"John","TelNo":"1234567890","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}`)},
		//test case 3: new firstname in existing Surname
		{"/", 200, []byte(`{"Surname":"Smith","Entries":[{"Firstname":"Jane","TelNo":"1234567890","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}`)},
		//test case 4: overwrite firstname in existing Surname
		{"/", 200, []byte(`{"Surname":"Smith","Entries":[{"Firstname":"John","TelNo":"0987654321","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}`)},
	}

	for i := range inputHTTP {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("PUT", inputHTTP[i].urlStr, bytes.NewBuffer(inputHTTP[i].JsonStr))

		router.ServeHTTP(w, req)
		fmt.Println(w.Code)
		if w.Code != inputHTTP[i].expectedWCode {
			t.Error("PutHandler test case", i, "returned", w.Code, "instead of", inputHTTP[i].expectedWCode)
		}
	}

}

func TestListHandler(t *testing.T) {

	router := httprouter.New()
	router.GET("/", listHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Println(w.Code)
	if w.Code != 200 {
		t.Error("listHandler returned", w.Code, "instead of 200")
	}

	bodyByte, _ := ioutil.ReadAll(w.Body)
	if string(bodyByte) != `{"Phone Book":[{"Surname":"Smith","Entries":[{"FirstName":"John","TelNo":"0987654321","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"},{"FirstName":"Jane","TelNo":"1234567890","Line1":"1a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}]}` {
		t.Error("listHandler returned body JSON is incorrect")
	}
}

type MockURL struct {
	urlStr        string
	expectedWCode int
}

func TestGetHandler(t *testing.T) {

	router := httprouter.New()
	router.GET("/:surname", getEntryHandler)
	router.GET("/:surname/:firstname", getEntryHandler)

	inputHTTP := [4]MockURL{
		//test case 0: improper formatted JSON
		{"/Dobbs", 404},
		//test case 1: more than 1 element being input in "Entries"
		{"/Smith", 200},
		//test case 2: new Surname
		{"/Smith/John", 200},
		//test case 3: new firstname in existing Surname
		{"/Smith/Roberto", 404},
	}

	for i := range inputHTTP {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", inputHTTP[i].urlStr, nil)

		router.ServeHTTP(w, req)
		fmt.Println(w.Code)
		if w.Code != inputHTTP[i].expectedWCode {
			t.Error("PutHandler test case", i, "returned", w.Code, "instead of", inputHTTP[i].expectedWCode)
		}
	}
}

func TestDeleteHandler(t *testing.T) {
	router := httprouter.New()
	router.DELETE("/:surname", delEntryHandler)
	router.DELETE("/:surname/:firstname", delEntryHandler)

	inputHTTP := [4]MockURL{
		//test case 0: Delete fistname granularity
		{"/Smith/John", 200},
		//test case 1: Delete fistname (which doesn't exists) granularity
		{"/Smith/Roberto", 200},
		//test case 2: Delete Surname (which doesn't exists)
		{"/Dobbs", 200},
		//test case 3: Delete Surname
		{"/Smith", 200},
	}

	for i := range inputHTTP {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", inputHTTP[i].urlStr, nil)

		router.ServeHTTP(w, req)
		fmt.Println(w.Code)
		if w.Code != inputHTTP[i].expectedWCode {
			t.Error("PutHandler test case", i, "returned", w.Code, "instead of", inputHTTP[i].expectedWCode)
		}
	}

	//Get from bucket
	var v []byte
	BoltClient.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		v = b.Get([]byte("Smith"))
		return nil
	})
	if v != nil {
		t.Error("Smith test entry was not deleted")
	}

	BoltClient.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("test"))
		return nil
	})
}
