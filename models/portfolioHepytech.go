package models

import "time"

type (
	PortfolioHepytech struct {
		Id        		uint      `gorm:"primary_key" json:"id"`
		CategoryID      uint      `gorm:"column:category_id" json:"id_category"`
		Name      		string    `json:"name"`
		Description     string    `gorm:"text" json:"description"`
		Image      		string    `gorm:"text" json:"image"`
		Desc_problem    string    `gorm:"text" json:"desc_problem"`
		Desc_solution   string    `gorm:"text" json:"desc_solution"`
		Slug      		string    `gorm:"unique" json:"slug"`
		Flag      		uint      `json:"flag"`
		CreatedAt 		time.Time `json:"created_at"`
		UpdatedAt 		time.Time `json:"updated_at"`
		Category     	Category  `gorm:"foreignKey:CategoryID" json:""`
		Impact 	        []Impact  `json:"-"`
	}
)