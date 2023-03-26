package model

type Response struct {
	Elapsed  int64      `json:"elapsed"`
	Average  float64    `json:"average"`
	Products []*Product `json:"products"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
