package main

type Material struct {
	ID     string `gorm:"primaryKey" json:"id"`
	TeamID string `json:"teamId"`
	Url    string `json:"url"`
}

func GetMaterials(teamId string) ([]Material, error) {
	var materials []Material
	result := DB.Where("team_id = ?", teamId).Find(&materials)
	return materials, result.Error
}
