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

type Keysrsas struct {
	Keyname    string `gorm:"type: varchar(100); unique_index"`
	Privatekey string `gorm:"type: longtext"`
	Publickey  string `gorm:"type: longtext"`
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/keys", allKeys).Methods("GET")
	myRouter.HandleFunc("/keys/{id}", findById).Methods("GET")
	myRouter.HandleFunc("/decrypt/{id}", decryptPrivateKey).Methods("GET")
	myRouter.HandleFunc("/new-key", createNewKeysrsas).Methods("POST")
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

func findById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Fprintf(w, "Key: "+key)
	keysrsas := []Keysrsas{}
	db.First(&keysrsas, "keyname = ?", key)
	if len(keysrsas) == 0 {

		json.NewEncoder(w).Encode("Don't Found")
	} else {
		json.NewEncoder(w).Encode(keysrsas)
	}
}

func decryptPrivateKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Fprintf(w, "Key: "+key+" ")
	keysrsas := []Keysrsas{}
	db.First(&keysrsas, "keyname = ?", key)

	if len(keysrsas) == 0 {

		json.NewEncoder(w).Encode("Don't Found")
	} else {
		privatekey := keysrsas[0].Privatekey
		decrypt := decryptAes(privatekey)
		keysrsas[0].Privatekey = decrypt
		json.NewEncoder(w).Encode(keysrsas)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func createNewKeysrsas(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var keysrsa Keysrsas

	json.Unmarshal(reqBody, &keysrsa)
	publ, priv := keys()
	keysrsa.Privatekey = priv
	keysrsa.Publickey = publ
	db.Create(&keysrsa)
	fmt.Println("Endpoint Hit: Creating New key")
	json.NewEncoder(w).Encode(keysrsa)
}
