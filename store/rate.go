package store

const (
	tinkoffGroup = "SavingAccountTransfers"
	EUR          = "EUR"
	RUB          = "RUB"
	USD          = "USD"
)

type TinkoffRate struct {
	ResultCode string `json:"resultCode"`
	Payload    struct {
		LastUpdate struct {
			Milliseconds int `json:"milliseconds"`
		} `json:"lastUpdate"`
		Rates []struct {
			Category     string `json:"category"`
			FromCurrency struct {
				Code    int    `json:"code"`
				Name    string `json:"name"`
				StrCode string `json:"strCode"`
			} `json:"fromCurrency"`
			ToCurrency struct {
				Code    int    `json:"code"`
				Name    string `json:"name"`
				StrCode string `json:"strCode"`
			} `json:"toCurrency"`
			Buy  float32 `json:"buy,omitempty"`
			Sell float32 `json:"sell,omitempty"`
		} `json:"rates"`
	} `json:"payload"`
	TrackingID string `json:"trackingId"`
}

type SberRate struct {
	Base struct {
		Num840 struct {
			Num0 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int     `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"0"`
			Num1000 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int     `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"1000"`
		} `json:"840"`
		Num978 struct {
			Num0 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int     `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"0"`
			Num1000 struct {
				IsoCur          string  `json:"isoCur"`
				CurrencyName    string  `json:"currencyName"`
				CurrencyNameEng string  `json:"currencyNameEng"`
				RateType        string  `json:"rateType"`
				CategoryCode    string  `json:"categoryCode"`
				Scale           int     `json:"scale"`
				BuyValue        float32 `json:"buyValue"`
				SellValue       float32 `json:"sellValue"`
				ActiveFrom      int     `json:"activeFrom"`
				BuyValuePrev    float32 `json:"buyValuePrev"`
				SellValuePrev   float32 `json:"sellValuePrev"`
				AmountFrom      int     `json:"amountFrom"`
				AmountTo        float64 `json:"amountTo"`
			} `json:"1000"`
		} `json:"978"`
	} `json:"base"`
}

type Rate struct {
	Buy          float32 `json:"buy,omitempty"`
	Sell         float32 `json:"sell,omitempty"`
	ToCurrency   string  `json:"toCurrency"`
	FromCurrency string  `json:"fromCurrency"`
	LastUpdate   int     `json:"lastUpdate"`
}

func MakeFromTinkoff(tinkoffRate *TinkoffRate) *Rate {
	rates := tinkoffRate.Payload.Rates
	for _, r := range rates {
		if r.Category == tinkoffGroup {
			return &Rate{
				Buy:          r.Buy,
				Sell:         r.Sell,
				ToCurrency:   r.ToCurrency.Name,
				FromCurrency: r.FromCurrency.Name,
				LastUpdate:   tinkoffRate.Payload.LastUpdate.Milliseconds,
			}
		}
	}
	return nil
}

func MakeFromSber(sberRate *SberRate, amount int, currency string) *Rate {
	rate := &Rate{}
	rate.ToCurrency = RUB
	rate.FromCurrency = currency
	switch currency {
	case EUR:
		if amount >= 1000 {
			rate.Sell = sberRate.Base.Num978.Num1000.SellValue
			rate.Buy = sberRate.Base.Num978.Num1000.BuyValue
		} else {
			rate.Sell = sberRate.Base.Num978.Num0.SellValue
			rate.Buy = sberRate.Base.Num978.Num0.BuyValue
		}
	case USD:
		if amount >= 1000 {
			rate.Sell = sberRate.Base.Num840.Num1000.SellValue
			rate.Buy = sberRate.Base.Num840.Num1000.BuyValue
		} else {
			rate.Sell = sberRate.Base.Num840.Num0.SellValue
			rate.Buy = sberRate.Base.Num840.Num0.BuyValue
		}
	}
	return rate
}
