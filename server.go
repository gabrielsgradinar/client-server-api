package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main(){

	_, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// db.AutoMigrate(&Cotacao{})


	cotacao, err := getCotacao()
	if err != nil {
		panic(err)
	}

	// db.Create(&cotacao)

	fmt.Println(cotacao)
}


func getCotacao() (*Cotacao, error){
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil{
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var filterJson map[string]interface{}
	json.Unmarshal([]byte(res), &filterJson)

	jsonStr, err := json.Marshal(filterJson["USDBRL"])
    if err != nil {
        return nil, err
    }

	var data Cotacao
	err = json.Unmarshal(jsonStr, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}