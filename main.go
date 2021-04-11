package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/i-redbyte/bank-aggregator/rest"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	fmt.Println("OK...")
	router.HandleFunc("/api/v1/rates", rest.AllRates).Methods("GET")
	router.HandleFunc("/api/v1/rate", rest.RateOwner).Methods("GET")
	err := http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}
