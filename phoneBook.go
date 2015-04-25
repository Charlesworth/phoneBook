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

	//get the body of the HTTP request
	rbodyByte, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	//Unmarshall into SurnameStruct var
	var rBodyJSON SurnameStruct
	err := json.Unmarshal(rbodyByte, &rBodyJSON)

	if err != nil {
		//if doesn't unmarshal then send back the error
		fmt.Fprint(w, "failed to marshal JSON: ", err)
		w.WriteHeader(400) //http: multiple response.WriteHeader calls error

	} else {

		//check if that surname is present in the phoneBook entries
		BoltClient.Mutex.RLock()
		var v []byte
		BoltClient.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("phoneBook"))
			v = b.Get([]byte(rBodyJSON.Surname))
			return nil
		})
		BoltClient.Mutex.RUnlock()

		if v == nil {
			//if v == nil then theres no entry for that surname
			//make a new entry
			BoltClient.Mutex.Lock()
			err = BoltClient.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("phoneBook"))
				err := b.Put([]byte(rBodyJSON.Surname), []byte(rbodyByte))
				return err
			})
			BoltClient.Mutex.Unlock()

			if err != nil {
				log.Print(err)
			}

		} else {
			//v != nil, so there is already some content under that surname
			//unmarshal boltDB entry
			var boltDBJSON SurnameStruct
			json.Unmarshal(v, &boltDBJSON)

			//Now check if the first name already is present
			//check each of the first names in the entries slice and replace if present
			newFirstName := true
			for i := range boltDBJSON.Entries {
				if boltDBJSON.Entries[i].FirstName == rBodyJSON.Entries[0].FirstName {
					boltDBJSON.Entries[i] = rBodyJSON.Entries[0]
					newFirstName = false
					break
				}
			}

			//else if the first name is new, append the entry to the slice
			if newFirstName {
				boltDBJSON.Entries = append(boltDBJSON.Entries, rBodyJSON.Entries...)
			}

			//then marshal back to []byte ready for BoltDB
			newJSON, err := json.Marshal(boltDBJSON)
			if err != nil {
				fmt.Println(err)
			}

			//rewrite the entry to bolt
			BoltClient.Mutex.Lock()
			err = BoltClient.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("phoneBook"))
				err := b.Put([]byte(boltDBJSON.Surname), newJSON)
				return err
			})
			BoltClient.Mutex.Unlock()

			if err != nil {
				log.Print(err)
			}
		}
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
