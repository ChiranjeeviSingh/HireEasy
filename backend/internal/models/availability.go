package models

import ( 
    "time"
    )

type Availability struct {
	ID        int       `json:"id,omitempty" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Date      string    `json:"date" db:"date" binding:"required"`
	FromTime  string `json:"from_time" db:"from_time" binding:"required"`
	ToTime    string `json:"to_time" db:"to_time" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}