// Written by Charles Cochrane, 2015
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"runtime"
	//"sync"
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
	test()

	//start the server
	log.Println("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", NewRouter()))
}

func listHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("list")
	//for i ; range listOfSurnames{
	//entry = get from bolt[i]
	//perhaps unmarsall
	//fmt.fprintln(entry)
	//}
}

func searchHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println(params.ByName("surname"))

}

func getEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("get", params.ByName("surname"))
	// dataMutex.Lock()
	// value, keyExists := Names.m[params.ByName("surname")]
	// if keyExists {
	// 	//value = unmarshall(value)
	// 	fmt.Fprintln(w, value)
	// } else {
	// 	//print result "no such entry"
	// 	fmt.Fprintln(w, "Entry not found")
	// 	w.WriteHeader(404)
	// }
	// dataMutex.Unlock()
}

func putEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("put", params.ByName("surname"))

}

func delEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println(params.ByName("surname"))

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
	//set program to use all available proccessors
	runtime.GOMAXPROCS(2)
}

func test() {

	m1 := Entry{
		FirstName:   "Ted",
		TelNo:       "test number",
		Line1:       "string",
		Line2:       "string",
		TownCity:    "string",
		CountyState: "string",
		Country:     "string",
		ZipPostal:   "string",
	}

	m1json, err := json.Marshal(m1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(m1json))

	var unmarsh Entry
	err = json.Unmarshal(m1json, &unmarsh)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("unmarshaled ", string(unmarsh))
}
