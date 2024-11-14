package models

import (
	"time"
)

type (
	Form struct {
		Id        			uint      `gorm:"primary_key" json:"id"`
		Name      			string    `json:"name"`
		Email      			string    `json:"email"`
		Message     		string    `gorm:"text" json:"message"`
		Flag     			uint      `json:"flag"`
		CreatedAt 			time.Time `json:"created_at"`
		UpdatedAt 			time.Time `json:"updated_at"`
	}
)