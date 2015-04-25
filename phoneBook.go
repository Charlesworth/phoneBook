// Written by Charles Cochrane, 2015
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

type SurnameStruct struct {
	Surname string        `json:"Surname"`
	Entries []EntryStruct `json:"Entries"`
}

type EntryStruct struct {
	FirstName   string `json:"FirstName"`
	TelNo       string `json:"TelNo"`
	Line1       string `json:"Line1"`
	Line2       string `json:"Line2"`
	TownCity    string `json:"TownCity"`
	CountyState string `json:"CountyState"`
	Country     string `json:"Country"`
	ZipPostal   string `json:"ZipPostal"`
}

var BoltClient struct {
	DB    *bolt.DB
	Mutex *sync.RWMutex
}

func main() {

	SetProc()
	NewBoltClient()
	defer BoltClient.DB.Close()

	//start the server
	fmt.Println("PhoneBook HTTP service")
	log.Println("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", NewRouter()))
}

func listHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("list")

	fmt.Fprint(w, `{"Phone Book":`)
	BoltClient.Mutex.RLock()

	//list all entries in the bucket
	BoltClient.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("phoneBook"))
		b.ForEach(func(k, v []byte) error {
			fmt.Fprintf(w, "%s,", v) //"%s\n", v)
			return nil
		})
		return nil
	})

	BoltClient.Mutex.RUnlock()
	fmt.Fprint(w, `}`)
}

//func searchHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	log.Println("search", params.ByName("surname"))
//}

func getEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("get", params.ByName("surname"))

	BoltClient.Mutex.RLock()

	//Get from bucket
	BoltClient.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("phoneBook"))
		v := b.Get([]byte(params.ByName("surname")))
		if v == nil {
			w.WriteHeader(404)
		} else {
			fmt.Fprintf(w, "%s\n", v)
		}
		return nil
	})

	BoltClient.Mutex.RUnlock()

}

func putEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("put", params.ByName("surname"))

	//get the content of the HTTP request

	bodyByte, _ := ioutil.ReadAll(r.Body)
	//body := string(bodyByte)
	r.Body.Close()
	//log.Println("Content: " + body)	test line

	//marshal Body into Json

	//	p, err := json.Marshal(body)
	//	if err != nil {
	//		fmt.Println(err)
	//	}

	//	fmt.Println(string(p))

	//Unmarshall into SurnameStruct var
	var unmarshal SurnameStruct
	err := json.Unmarshal(bodyByte, &unmarshal)

	if err != nil {
		//if doesn't unmarshal then send back bad request
		fmt.Fprint(w, "failed to marshal JSON: ", err)
		w.WriteHeader(400) //http: multiple response.WriteHeader calls error

	} else {

		//check if that surname is present
		BoltClient.Mutex.RLock()
		var v []byte
		BoltClient.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("phoneBook"))
			v = b.Get([]byte(unmarshal.Surname))
			return nil
		})
		BoltClient.Mutex.RUnlock()

		if v == nil { //doesn't seem to work after a delete
			fmt.Println("entry does NOT exist")
			//just write the thing
			BoltClient.Mutex.Lock()

			err = BoltClient.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("phoneBook"))
				err := b.Put([]byte(unmarshal.Surname), []byte(bodyByte))
				return err
			})

			BoltClient.Mutex.Unlock()

			if err != nil {
				log.Print(err)
			}

		} else {
			//insert into current
			fmt.Println("does exist")

			//unmarshal boltDB copy
			var unmarshal2 SurnameStruct
			json.Unmarshal(v, &unmarshal2)

			//delete if there is already one with this first name
			newName := true
			for i := range unmarshal2.Entries {
				if unmarshal2.Entries[i].FirstName == unmarshal.Entries[0].FirstName {
					unmarshal2.Entries[i] = unmarshal.Entries[0]
					newName = false
					break
				}
			}

			if newName {
				//append []entries with new entry
				unmarshal2.Entries = append(unmarshal2.Entries, unmarshal.Entries...)
			}

			//then marshal
			m1json, err := json.Marshal(unmarshal2)
			if err != nil {
				fmt.Println(err)
			}

			//rewrite to bolt
			BoltClient.Mutex.Lock()

			err = BoltClient.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("phoneBook"))

				err := b.Put([]byte(unmarshal2.Surname), m1json)
				return err
			})

			BoltClient.Mutex.Unlock()

			if err != nil {
				log.Print(err)
			}

		}

		//if it is then retrieve and add record
		//else

		//		w.WriteHeader(200)
	}
}

func delEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("del", params.ByName("surname"), params.ByName("firstname"))

	//get surname
	//if not present return 404
	//else check if firstname is present
	//if not present return 404
	//else remove that part of the JSON
	BoltClient.Mutex.Lock()

	//Delete from bucket
	err := BoltClient.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("phoneBook"))
		err := b.Delete([]byte(params.ByName("surname")))
		return err
	})

	BoltClient.Mutex.Unlock()

	if err != nil {
		log.Print(err)
	}

}

//NewRouter returns a httprouter.Router complete with the routes
func NewRouter() *httprouter.Router {

	router := httprouter.New()
	router.GET("/list", listHandler)
	//router.GET("/search/:surname", searchHandler) //need to fix this, to /search?surname=bob
	router.GET("/entry/:surname", getEntryHandler) //is this the search
	router.PUT("/entry/:surname", putEntryHandler)
	router.DELETE("/entry/:surname/:firstname", delEntryHandler)

	return router
}

//SetProc sets the program to use 2 proccessor cores
func SetProc() {
	runtime.GOMAXPROCS(2)
}

//NewBoltClient produces a BoltClient and Mutex lock and assigns them
//to the global BoltClient variable
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
}
