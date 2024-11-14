package models

import "time"

type (
	Category struct {
		ID        uint      			`gorm:"primary_key" json:"id"`
		Name      string    			`json:"name"`	
		Flag      uint      			`json:"flag"`
		CreatedAt time.Time 			`json:"created_at"`
		UpdatedAt time.Time 			`json:"updated_at"`
		Career 	  []Career				`json:"-"`
		Portfolio []PortfolioHepytech 	`json:"-"`
	}
)