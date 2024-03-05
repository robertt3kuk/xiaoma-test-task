package model

import "time"

//swagger:model
type Customer struct {
	ID        int        `json:"id"`
	Name      string     `json:"customer_name"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
