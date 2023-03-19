package main

import (
	"encoding/json"
	"fmt"
	"giftcard/connect"
	member "giftcard/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var version = "v1"
var urlhder = "/" + version

var ipServer = ""
var portServer = ""

// const (
// 	// Secret key used to sign JWTs
// 	secret = "my-secret-key"
// )

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	// Add the auth middleware to the router
	//router.Use(authMiddleware)
	router.HandleFunc(urlhder, home)
	router.HandleFunc(urlhder+"/cardid/{id}", getUserCardId)
	router.HandleFunc(urlhder+"/membid/{id}", getUserMembId)
	router.HandleFunc(urlhder+"/membidname/{id}", getUserMembIdLoan)
	router.HandleFunc(urlhder+"/membname/{name}", getUserMembName)
	router.HandleFunc(urlhder+"/updatememb", updateUserMemb).Methods("POST")
	fmt.Println("#-------------------------------------------------------------------------#")
	fmt.Println("#           Server start on  " + ipServer + ":" + portServer)
	fmt.Println("#-------------------------------------------------------------------------#")
	log.Fatal(http.ListenAndServe(ipServer+":"+portServer, router))

}

func main() {

	dataBaseConnection()
	handleRequests()
}

func dataBaseConnection() {
	connect.Connect()

}

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(" Hello, GIFT VOTE...! ")
}

func getUserCardId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	var u = connect.GetUserCardid(id)

	fmt.Println(u)

	json.NewEncoder(w).Encode(u)

}

func getUserMembId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	var u = connect.GetUserMembId(id)

	fmt.Println(u)

	json.NewEncoder(w).Encode(u)
}

func getUserMembIdLoan(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	var u = connect.GetUserMembIdLoan(id)

	fmt.Println(u)

	json.NewEncoder(w).Encode(u)
}

func getUserMembName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	var u = connect.GetUserMembName(name)

	fmt.Println(u)

	json.NewEncoder(w).Encode(u)

}

func updateUserMemb(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var memb []member.MemberArray
	json.Unmarshal(reqBody, &memb)
	err := connect.UpdateUser(memb)
	//member.MemberArray  = append(member.MemberArray , article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Err"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	//json.NewEncoder(w).Encode(memb)

}
