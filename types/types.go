package types

var URL_API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type CotacaoDataDTO struct {
	Code       string `json:"code"`
	CodeIN     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	BID        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
