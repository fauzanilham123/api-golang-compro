package models

import (
	"time"
)

type (
	Home struct {
		Id        						uint      `gorm:"primary_key" json:"id"`
		Logo      						string    `gorm:"text" json:"logo"`
		Background_image_section_1      string    `gorm:"text" json:"background_image_section_1"`
		Title_section_1      			string    `json:"title_section_1"`
		Description_section_1      		string    `gorm:"text" json:"description_section_1"`
		Button_section_1     			string    `json:"button_section_1"`
		Sub_title_section_2      		string    `json:"sub_title_section_2"`
		Title_section_2      			string    `json:"title_section_2"`
		Description_section_2      		string    `gorm:"text" json:"description_section_2"`
		Button_section_2      			string    `json:"button_section_2"`
		Image_section_2      			string    `gorm:"text" json:"image_section_2"`
		Sub_title_section_3      		string    `json:"sub_title_section_3"`
		Title_section_3      			string    `json:"title_section_3"`
		Sub_title_section_4      		string    `json:"sub_title_section_4"`
		Title_section_4      			string    `json:"title_section_4"`
		Description_contact_us      	string    `gorm:"text" json:"description_contact_us"`
		Button_contact_us      			string    `json:"button_contact_us"`
		Sub_title_section_5      		string    `json:"sub_title_section_5"`
		Title_section_5      			string    `json:"title_section_5"`
		Button_section_5      			string    `json:"button_section_5"`
		Link_facebook      			    string    `gorm:"text" json:"link_facebook"`
		Link_linkedln      			    string    `gorm:"text" json:"link_linkedln"`
		Link_instagram      			string    `gorm:"text" json:"link_instagram"`
		Flag      						uint      `json:"flag"`
		CreatedAt 						time.Time `json:"created_at"`
		UpdatedAt 						time.Time `json:"updated_at"`
	}
)