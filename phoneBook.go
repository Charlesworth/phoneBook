// Written by Charles Cochrane, 2015
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"runtime"
	"sync"
)

var Names = struct {
	sync.RWMutex
	m map[string][]string
}{m: make(map[string][]string)}

type Entry struct {
	TelNo   string
	Address AddressStruct
}

type AddressStruct struct {
	Line1       string
	Line2       string
	TownCity    string
	CountyState string
	Country     string
	ZipPostal   string
}

func main() {

	SetProc()

	//start the server
	log.Println("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", NewRouter()))
}

func ListHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	log.Println(params.ByName("surname"))
}

func putEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println(params.ByName("surname"))
}

func delEntryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println(params.ByName("surname"))
}

func NewRouter() *httprouter.Router {

	router := httprouter.New()
	router.GET("/list", ListHandler)
	router.GET("/search/:surname", searchHandler) //need to fix this, to /search?surname=bob
	router.GET("/entry/:surname", getEntryHandler)
	router.PUT("/entry/:surname", putEntryHandler)
	router.DELETE("/entry/:surname", delEntryHandler)

	return router
}

func SetProc() {
	//set program to use all available proccessors
	//procNo := runtime.NumCPU()
	runtime.GOMAXPROCS(2)
}

func test() {

}
