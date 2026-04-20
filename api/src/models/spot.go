package models

import "gorm.io/gorm"

type Spot struct {
	gorm.Model
	UserId uint `json:"user_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Description string `json:"description"`
	OpenFrom string `json:"open_from"`
	OpenTo string `json:"open_to"`
}

func (spot *Spot) Save() (*Spot, error) {
	err := Database.Create(spot).Error
	if err != nil {
		return &Spot{}, err
	}
	return spot, nil
}

func FetchAllSpots() (*[]Spot, error) {
	var spots []Spot
	err := Database.Find(&spots).Error
	if err != nil {
		return &[]Spot{}, err
	}
	return &spots, nil
}

func FindSpot(id string) (*Spot, error) {
	var spot Spot
	err := Database.Where("id = ?", id).First(&spot).Error
	if err != nil {
		return &Spot{}, err
	}
	return &spot, nil
}