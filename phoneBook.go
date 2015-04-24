// Written by Charles Cochrane, 2015
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

type Entry struct {
	FirstName   string `json:"FirstName"`
	TelNo       string `json:"TelNo"`
	Line1       string `json:"Line1"`
	Line2       string `json:"Line2"`
	TownCity    string `json:"TownCity"`
	CountyState string `json:"CountyState"`
	Country     string `json:"Country"`
	ZipPostal   string `json:"ZipPostal"`
}

func main() {

	SetProc()
	NewBoltClient()
	defer BoltClient.DB.Close()

	//start the server
	log.Println("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", NewRouter()))
}

func listHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("list")
}

func searchHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("search", params.ByName("surname"))

}

func getEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("get", params.ByName("surname"))
}

func putEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("put", params.ByName("surname"))

}

func delEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("del", params.ByName("surname"))

}

func NewRouter() *httprouter.Router {

	router := httprouter.New()
	router.GET("/list", listHandler)
	router.GET("/search/:surname", searchHandler) //need to fix this, to /search?surname=bob
	router.GET("/entry/:surname", getEntryHandler)
	router.PUT("/entry/:surname", putEntryHandler)
	router.DELETE("/entry/:surname", delEntryHandler)

	return router
}

func SetProc() {
	//set program to use 2 proccessor cores
	runtime.GOMAXPROCS(2)
}

//Object variables
var BoltClient struct {
	DB    *bolt.DB
	Mutex *sync.RWMutex
}

//constructor function
func NewBoltClient() { //*BoltClient {
	//Open DB
	BoltDB, err := bolt.Open("phoneBook.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	//open (or create if not present) bucket
	err = BoltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("phoneBook"))

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	//make a read/write mutex
	mutex := &sync.RWMutex{}

	//assign the pointers to the mutex and DB to the Global
	//BoltClient var
	BoltClient.DB = BoltDB
	BoltClient.Mutex = mutex

	//	return &BoltClient{
	//		BoltDB,
	//		mutex,
	//	}
}
