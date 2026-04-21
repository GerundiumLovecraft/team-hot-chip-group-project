package models

import "gorm.io/gorm"

type Spot struct {
	gorm.Model
	UserId      uint           `json:"user_id"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	Description string         `json:"description"`
	OpenFrom    string         `json:"open_from"`
	OpenTo      string         `json:"open_to"`
	Features    []SpotsToFeats `json:"features" gorm:"foreignKey:SpotId;references:ID;constraint:OnDelete:CASCADE"`
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
	err := Database.Preload("Features.Feature").Find(&spots).Error
	if err != nil {
		return &[]Spot{}, err
	}
	return &spots, nil
}

func FindSpot(id string) (*Spot, error) {
	var spot Spot
	err := Database.Preload("Features.Feature").First(&spot, id).Error
	if err != nil {
		return &Spot{}, err
	}
	return &spot, nil
}

func FilterSpotsByFeature(id uint, value *int8) (*[]Spot, error) {
	var spots []Spot

	// initiate construction of the query
	query := Database.Joins("JOIN spots_to_feats ON spots_to_feats.spot_id = spots.id").
		Where("spots_to_feats.feat_id = ?", id)

	// add additional filter if value is not nil
	if value != nil {
		query = query.Where("spots_to_feats.value = ?", value)
	}

	// send query to DB
	err := query.Preload("Features.Feature").Find(&spots).Error

	if err != nil {
		return &[]Spot{}, err
	}

	return &spots, nil
}
