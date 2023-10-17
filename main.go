package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/Mule/getRepo", GetRepository).Methods(http.MethodGet).Name("GetAllBook")

	fmt.Println("Server is getting started...")

	fmt.Println("Listening at port 4000..")

	log.Fatal(http.ListenAndServe("localhost:8080", router))

	http.Handle("/", router)

}
