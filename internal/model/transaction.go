package model

import "time"

//swagger:model
type Transaction struct {
	ID         int        `json:"id"`
	CustomerID int        `json:"customer_id"`
	ItemID     int        `json:"item_id"`
	Qty        int        `json:"qty"`
	Price      float64    `json:"price"`
	Amount     float64    `json:"amount"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

//swagger:model
type TransactionView struct {
	ID           int        `json:"id"`
	CustomerID   int        `json:"customer_id"`
	CustomerName string     `json:"customer_name"`
	ItemID       int        `json:"item_id"`
	ItemName     string     `json:"item_name"`
	Qty          int        `json:"qty"`
	Price        float64    `json:"price"`
	Amount       float64    `json:"amount"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
