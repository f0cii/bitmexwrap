package bitmexwrap

type Symbol struct {
	Name          string
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
}
