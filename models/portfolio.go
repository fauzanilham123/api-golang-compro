package models

import "time"

type (
	Portfolio struct {
		Id               uint      		  `gorm:"primary_key" json:"id"`
		Image       	 string    		  `gorm:"text" json:"image" form:"image"`
		Title       	 string    		  `json:"title" form:"title"`
		Category_homeID  uint      		  `json:"id_category" form:"id_category"`
		Flag             uint      		  `json:"flag"`
		CreatedAt        time.Time 		  `json:"created_at"`
		UpdatedAt   	 time.Time 		  `json:"updated_at"`
		Category    	 Category_home    `gorm:"foreignKey:Category_homeID" json:""`
	}
)