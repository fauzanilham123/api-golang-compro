package models

import "time"

type (
	Service struct {
		Id        		uint      `gorm:"primary_key" json:"id"`
		Icon      		string    `gorm:"text" json:"icon"`
		Title     		string    `json:"title"`
		Description     string    `json:"description"`
		Flag      		uint      `json:"flag"`
		CreatedAt 		time.Time `json:"created_at"`
		UpdatedAt 		time.Time `json:"updated_at"`
	}
)