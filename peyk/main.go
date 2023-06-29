package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:my-secret-pw@tcp(localhost:3306)/bepa")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := mux.NewRouter()
	// Add the cors.Default() middleware
	handler := cors.Default().Handler(router)
	router.HandleFunc("/subscription", addSubscription).Methods("POST")
	router.HandleFunc("/history/{coin_name}", getHistory).Methods("GET")
	fmt.Println(http.ListenAndServe(":8081", handler))
}

type Price struct {
	CoinName  string `json:"coin_name"`
	Timestamp string `json:"timestamp"`
	Price     string `json:"price"`
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coin_name := params["coin_name"]

	rows, err := db.Query("SELECT * FROM prices WHERE coin_name = ?", coin_name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prices []Price
	for rows.Next() {
		var p Price
		err := rows.Scan(&p.CoinName, &p.Timestamp, &p.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prices = append(prices, p)
	}
	json.NewEncoder(w).Encode(prices)
}

type AlertSubscription struct {
	UserEmail            string  `json:"user_email"`
	CoinName             string  `json:"coin_name"`
	DifferencePercentage float64 `json:"difference_percentage"`
}

func addSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription AlertSubscription
	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare("INSERT INTO alertSubscription(user_email, coin_name, difference_percentage) VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(subscription.UserEmail, subscription.CoinName, subscription.DifferencePercentage)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New subscription was added")
}
