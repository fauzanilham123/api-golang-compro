package models

import "time"

type (
	Impact struct {
		Id           			uint      `gorm:"primary_key" json:"id"`
		PortfolioHepytechId 	uint      `gorm:"column:portfolio_hepytech_id" json:"id_portfolio"`
		Impact_title 			string    `json:"impact_title"`
		Impact_desc  			string    `gorm:"text" json:"impact_desc"`
		Impact_icon  			string    `gorm:"text" json:"impact_icon"`
		Flag         			uint      `json:"flag"`
		CreatedAt    			time.Time `json:"created_at"`
		UpdatedAt    			time.Time `json:"updated_at"`
		Portfolio     			PortfolioHepytech  `gorm:"foreignKey:PortfolioHepytechId" json:""`
	}
)