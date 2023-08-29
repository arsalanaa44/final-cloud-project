package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/arsalanaa44/final-cloud-project/bepa/models"
	"github.com/arsalanaa44/final-cloud-project/bepa/internal/mail"
)

func main() {
	getData()
	handleSubscription()
}

func getData() {
	coins := getAvailableCoins()

	for _, coin := range coins {
		info := getCoinInfo(coin)
		saveCoinInfo(info)
	}
}

func getAvailableCoins() []string {
	resp, err := http.Get("http://localhost:8000/api/data")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var coins []string
	json.Unmarshal(body, &coins)

	return coins
}

func getCoinInfo(name string) models.Coin {
	resp, err := http.Get("http://localhost:8000/api/data/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var coin models.Coin
	json.Unmarshal(body, &coin)

	return coin
}

func saveCoinInfo(coin models.Coin) {
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(db:3306)/bepa")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO prices(coin_name, timestamp, price) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(coin.Name, coin.UpdatedAt, coin.Value)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSubscription() {

	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:3306)/bepa")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	latestPrices, err := mail.GetLatestPrices(db)
	if err != nil {
		log.Fatal(err)
	}

	alerts, err := mail.CheckAlerts(db, latestPrices)
	if err != nil {
		log.Fatal(err)
	}

	err = mail.SendEmails(alerts)
	if err != nil {
		log.Fatal(err)
	}
}
