/*
Package restapi Generate Key Public and Private

Lab Web Development Go

    Host: localhost:8083
    BasePath: /
    Version: 0.1.0
	Title: Lab Web Development Go
swagger:meta
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// Keysrsas
//
// Structure keys private and public
//
// swagger:model Keysrsas
type Keysrsas struct {
	// name of the keys
	// example: peterparker
	// required: true
	// keyname
	Keyname    string `gorm:"type: varchar(100); unique_index"`
	Privatekey string `gorm:"type: longtext"`
	Publickey  string `gorm:"type: longtext"`
}

// swagger:parameters keyname
type _ struct {
	// name of the keys
	// in:path
	// required: true
	Keyname string `json:"keyname"`
}

// swagger:parameters id
type _ struct {
	// name of the keys
	// in:path
	// required: true
	Keyname string `json:"id"`
}

// swagger:parameters keysrsas
type _ struct {
	// The body to create a thing
	// in:body
	// required: true
	Body string `json:"keyname"`
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/keys", allKeys).Methods("GET")
	myRouter.HandleFunc("/key/{keyname}", findById).Methods("GET")
	myRouter.HandleFunc("/key/decrypt/{id}", decryptPrivateKey).Methods("GET")
	myRouter.HandleFunc("/key", createNewKeysrsas).Methods("POST")
	myRouter.HandleFunc("/swagger", getSwagger).Methods("GET")
	log.Fatal(http.ListenAndServe(":8083", myRouter))
}

func main() {

	db, err = gorm.Open("mysql",
		"rw:password@/rsa?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("connection succedssed")
	}

	db.AutoMigrate(&Keysrsas{})

	handleRequests()

}
func allKeys(w http.ResponseWriter, r *http.Request) {

	keysrsas := []Keysrsas{}
	db.Find(&keysrsas)
	json.NewEncoder(w).Encode(keysrsas)

}

// swagger:route GET /key/{keyname} findById keyname
// Obtain  keys by keyname
// Responses:
// - 200: Keysrsas
// - 404: Keysrsas
func findById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	keysrsas := findKeyname(w, key)
	db.First(&keysrsas, "keyname = ?", key)
	if len(keysrsas) != 0 {
		json.NewEncoder(w).Encode(keysrsas)
	}
}

// swagger:route GET /key/decrypt/{id} decryptPrivateKey id
// Obtain keys and Decrypt private key by  keyname
// Responses:
// - 200: Keysrsas
// - 404: Keysrsas
func decryptPrivateKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	keysrsas := findKeyname(w, key)

	if len(keysrsas) != 0 {
		privatekey := keysrsas[0].Privatekey
		decryptsa := decryptAes(privatekey)
		keysrsas[0].Privatekey = (decryptsa)
		json.NewEncoder(w).Encode(keysrsas)
	}
}

func findKeyname(w http.ResponseWriter, keyname string) []Keysrsas {
	keysrsas := []Keysrsas{}
	db.First(&keysrsas, "keyname = ?", keyname)
	if len(keysrsas) == 0 {
		keysrsa := Keysrsas{}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(keysrsa)
	}
	return keysrsas
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome to HomePage!</h1>")
	fmt.Fprintf(w, "<a href=http://localhost:8083/keys>All Keys</a><br>")
	fmt.Fprintf(w, "<a href=http://localhost:8083/swagger> Swagger</a>")

}

func getSwagger(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("./swagger.json")

	check(err)
	fmt.Fprintf(w, string(dat))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// swagger:route POST /key createNewKey keysrsas
// Save public and private key  keyname
// Responses:
// - 201: Keysrsas
func createNewKeysrsas(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var keysrsa Keysrsas

	json.Unmarshal(reqBody, &keysrsa)
	publ, priv := keys()

	keysrsa.Privatekey = priv
	keysrsa.Publickey = publ
	keysrs := findKeyname(w, keysrsa.Keyname)
	if len(keysrs) != 0 {
		fmt.Println("Endpoint Hit: Key Exists")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Key Exists")
	} else {
		db.Create(&keysrsa)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(keysrsa)
	}
}
