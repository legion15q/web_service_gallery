//main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello")
	router := mux.NewRouter()
	router.HandleFunc("/", DoHealthCheck).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func DoHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, i'm a golang microservice")
	fmt.Println(r.URL)
	w.WriteHeader(http.StatusAccepted) //RETURN HTTP CODE 202
}
