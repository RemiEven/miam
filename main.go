package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type product struct {
	ID int `json:"id"`
}

func productHandler(responseWriter http.ResponseWriter, request *http.Request) {
	json.NewEncoder(responseWriter).Encode(product{
		ID: 1994,
	})
}

func main() {
	fmt.Println("Coucou")
	router := mux.NewRouter()
	router.HandleFunc("/product/1", productHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Coucou34")
}
