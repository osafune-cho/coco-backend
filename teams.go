package main

import "time"

type Team struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	CourseID  string `json:"courseId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetTeam(id string) (*Team, error) {
	var team Team
	result := DB.Where("id = ?", id).First(&team)
	return &team, result.Error
}
