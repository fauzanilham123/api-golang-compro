package models

import "time"

type (
	Navbar struct {
		Id          uint      `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		Link_button string    `gorm:"type:text" json:"link_button"`
		Flag        uint      `json:"flag"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
)