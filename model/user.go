package model

import "time"

//User struct declaration
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_time" gorm:"autoUpdateTime:milli"`
}
