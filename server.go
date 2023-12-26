package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const tableName = "cotacoes"

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

func (Cotacao) TableName() string {
	return tableName
}


func main(){
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.Table(tableName).AutoMigrate(&Cotacao{})

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		getCotacaoHandler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := getCotacaoData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db.Create(&cotacao)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	result, err := json.Marshal(cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func getCotacaoData() (*Cotacao, error){
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

	// armazena a resporta em um map
	var filterJson map[string]interface{}
	err = json.Unmarshal([]byte(res), &filterJson)
	if err != nil {
        return nil, err
    }

	// transformas os dados dentro da chave USDBRL para json novamente
	jsonStr, err := json.Marshal(filterJson["USDBRL"])
    if err != nil {
        return nil, err
    }

	// faz o unmarshal desses dados para a struct de Cotacao
	var data Cotacao
	err = json.Unmarshal(jsonStr, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}