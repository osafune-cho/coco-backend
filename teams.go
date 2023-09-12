package main

import "time"

type Team struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	CourseID  string `json:"courseId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
