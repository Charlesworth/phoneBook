package main

import (
	"github.com/julienschmidt/httprouter"
	//"net/http"
	//"net/http/httptest"
	//"github.com/Charlesworth/phoneBook"
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
	NewBoltClient()

}
