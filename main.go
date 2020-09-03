package main

import (
	"github.com/gorilla/mux"
	"github.com/ilya-sokolov/bank-aggregator/rest"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/rates", rest.AllRates).Methods("GET")
	router.HandleFunc("/api/v1/rate", rest.RateOwner).Methods("GET")
	err := http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
	//for {
	//	printAllRates()
	//	time.Sleep(10 * time.Second)
	//}
}

//func printAllRates() {
//	rate := getAllRates()
//	for _, r := range rate {
//		fmt.Printf("%s: %s -> %s BUY: %.2f  SELL: %.2f\n", r.Owner, r.FromCurrency, r.ToCurrency, r.Buy, r.Sell)
//	}
//	fmt.Println("<========================>")
//}
