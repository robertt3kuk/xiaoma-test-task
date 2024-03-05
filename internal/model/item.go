package model

import "time"

//swagger:model
type Item struct {
	ID        int        `json:"id"`
	ItemName  string     `json:"item_name"`
	Cost      float64    `json:"cost"`
	Price     float64    `json:"price"`
	Sort      int        `json:"sort"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
