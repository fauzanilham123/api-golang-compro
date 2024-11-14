package models

import (
	"time"
)

type (
	// AgeRatingCategory
	Career struct {
	ID          	uint      	`gorm:"primary_key" json:"id"`
	CategoryID 		uint    	`gorm:"column:category_id" json:"id_category"`
	PositionID 		uint		`gorm:"column:position_id" json:"id_position"`
	Name 			string 		`json:"name"`
	Description 	string 		`json:"description"`
	Required 		string 		`json:"required"`
	Flag			uint 		`json:"flag"`
	CreatedAt  		time.Time 	`json:"created_at"`
	UpdatedAt   	time.Time	`json:"updated_at"`
	Category     	Category    `gorm:"foreignKey:CategoryID" json:""`
	Position     	Position    `gorm:"foreignKey:PositionID" json:""`
	}
)