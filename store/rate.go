package store

const (
	tinkoffGroup = "SavingAccountTransfers"
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

type Rate struct {
	Buy          float32 `json:"buy,omitempty"`
	Sell         float32 `json:"sell,omitempty"`
	ToCurrency   string  `json:"toCurrency"`
	FromCurrency string  `json:"fromCurrency"`
	LastUpdate   int     `json:"lastUpdate"`
}

func (rate Rate) MakeFromTinkoff(tinkoffRate *TinkoffRate) *Rate {
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
