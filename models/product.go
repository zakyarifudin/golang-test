package models

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Barcode   string    `gorm:"type:varchar(100);unique;not null" json:"barcode"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	Image     string    `gorm:"type:varchar(255)" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
