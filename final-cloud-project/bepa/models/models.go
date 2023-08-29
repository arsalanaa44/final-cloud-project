package models


type Coin struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	UpdatedAt string  `json:"updated_at"`
}