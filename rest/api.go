package rest

import (
	"encoding/json"
	"fmt"
	"github.com/i-redbyte/bank-aggregator/store"
	"net/http"
	"sync"
	"time"
)

const (
	tinkoffApi = "https://api.tinkoff.ru/v1/currency_rates"
	sberApi    = "https://www.sberbank.ru/portalserver/proxy/?pipe=shortCachePipe&url=http%3A%2F%2Flocalhost%2Frates-web%2FrateService%2Frate%2Fcurrent%3FregionId%3D77%26rateCategory%3Dbase%26currencyCode%3D978%26currencyCode%3D840"
	alfaApi    = "https://alfabank.ru/ext-json/0.2/exchange/cash?offset=0&limit=1&mode=rest"
)

var client = &http.Client{Timeout: 10 * time.Second}

func getTinkoffRate(from, to string) (*store.Rate, error) {
	r, err := client.Get(tinkoffApi + "?from=" + from + "&to=" + to)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	tr := &store.TinkoffRate{}
	err = json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		return nil, err
	}
	return store.MakeFromTinkoff(tr), nil
}

func getSberRate(currency string, amount int) (*store.Rate, error) {
	r, err := client.Get(sberApi)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	sr := &store.SberRate{}
	err = json.NewDecoder(r.Body).Decode(sr)
	if err != nil {
		return nil, err
	}
	return store.MakeFromSber(sr, amount, currency), nil
}

func getAlfaRate(currency string) (*store.Rate, error) {
	r, err := client.Get(alfaApi)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	ar := &store.AlfaRate{}
	err = json.NewDecoder(r.Body).Decode(ar)
	if err != nil {
		return nil, err
	}
	return store.MakeFromAlfa(ar, currency), nil
}

func AllRates(w http.ResponseWriter, r *http.Request) {
	rates := getAllRates(store.USD)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
	for _, r := range rates {
		fmt.Printf("%s: %s -> %s BUY: %.2f  SELL: %.2f\n", r.Owner, r.FromCurrency, r.ToCurrency, r.Buy, r.Sell)
	}
}

func RateOwner(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	owner := query.Get("owner")
	currency := query.Get("currency")
	if currency == "" {
		currency = store.USD
	}
	rate := &store.Rate{}
	var err error
	switch owner {
	case "tinkoff":
		rate, err = getTinkoffRate(currency, store.RUB)
	case "sber":
		rate, err = getSberRate(currency, 100)
	case "alfa":
		rate, err = getAlfaRate(currency)
	}
	if err != nil {
		http.Error(w, "Field \"owner\" has empty.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

func getAllRates(currency string) (rates []store.Rate) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		rate, err := getTinkoffRate(currency, store.RUB)
		if err != nil {
			fmt.Println("TinkoffError ", err)
			wg.Done()
			return
		}
		rates = append(rates, *rate)
		wg.Done()
	}()
	go func() {
		rate, err := getSberRate(currency, 100)
		if err != nil {
			fmt.Println("SberError ", err)
			wg.Done()
			return
		}
		rates = append(rates, *rate)
		wg.Done()
	}()
	go func() {
		rate, err := getAlfaRate(currency)
		if err != nil {
			fmt.Println("AlfaError ", err)
			wg.Done()
			return
		}
		rates = append(rates, *rate)
		wg.Done()
	}()
	wg.Wait()
	return rates
}
