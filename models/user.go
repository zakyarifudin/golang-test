package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(100);unique" json:"username"`
	Password  string    `json:"-"`                            // "-" artinya password gak bakal muncul pas kita balikin data JSON
	Role      string    `gorm:"type:varchar(50)" json:"role"` // Misal: admin atau kasir
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
