package mail

import (
	"database/sql"
	"fmt"
	"net/smtp"
)

type Alert struct {
	Email      string
	CoinName   string
	Threshold  float64
	Difference float64
}

const (
	latestPricesSQL = `
        SELECT p1.coin_name, p1.price 
        FROM prices p1
        JOIN (
            SELECT coin_name, MAX(timestamp) as latest_time 
            FROM prices 
            GROUP BY coin_name
        ) p2 ON p1.coin_name = p2.coin_name AND p1.timestamp = p2.latest_time
    `

	alertsSQL = `
			SELECT a.user_email, a.coin_name, a.difference_percentage, p.price_difference
			FROM alertSubscription a
			JOIN (
				SELECT current.coin_name, 
					   (current.price - previous.price) / previous.price * 100 as price_difference
				FROM prices current
				JOIN prices previous ON current.coin_name = previous.coin_name
				WHERE current.timestamp = (
						SELECT MAX(timestamp) 
						FROM prices 
						WHERE coin_name = current.coin_name
					 )
				AND previous.timestamp = (
						SELECT MAX(timestamp) 
						FROM prices 
						WHERE coin_name = current.coin_name AND timestamp < current.timestamp
					 )
			) p ON a.coin_name = p.coin_name
			WHERE ABS(p.price_difference) > a.difference_percentage
		`
)

func SendEmails(alerts []Alert) error {
	for _, alert := range alerts {
		msg := fmt.Sprintf("Subject: Price Alert\n\nThe price of %s has changed by %.2f%%. Please check your currency.", alert.CoinName, alert.Difference)
		err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", "mehran.akmah@gmail.com", "oemcglbamvusrrht", "smtp.gmail.com"), "mehran.akmah@gmail.com", []string{alert.Email}, []byte(msg))
		if err != nil {
			return err
		}
		fmt.Println(msg + alert.Email)
	}
	return nil
}

func GetLatestPrices(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query(latestPricesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[string]float64)
	for rows.Next() {
		var coinName string
		var price float64
		if err := rows.Scan(&coinName, &price); err != nil {
			return nil, err
		}
		prices[coinName] = price
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return prices, nil
}
func CheckAlerts(db *sql.DB, latestPrices map[string]float64) ([]Alert, error) {
	rows, err := db.Query(alertsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts := make([]Alert, 0)
	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.Email, &alert.CoinName, &alert.Threshold, &alert.Difference); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return alerts, nil
}
