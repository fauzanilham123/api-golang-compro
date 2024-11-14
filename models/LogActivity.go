package models

import (
	"time"
)

type (
	LogActivity struct {
		Id        	uint      `gorm:"primary_key" json:"id"`
		UserID     	uint      `gorm:"column:user_id" json:"id_user"`
		Time      	time.Time `json:"time"`
		Activity    string    `json:"activity"`
		Method     	string    `json:"method"`
		User     	User      `gorm:"foreignKey:UserID" json:""`
	}
)