package models

import "time"

type (
	Logo struct {
		Id        uint      `gorm:"primary_key" json:"id"`
		Name      string    `json:"name"`
		Logo      string    `gorm:"text" json:"logo"`
		Flag      uint      `json:"flag"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)