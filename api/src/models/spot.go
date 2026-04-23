package models

import (
	"gorm.io/gorm"
)

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

type FeatureFilter struct {
	ID    uint  `json:"id"`
	Value *int8 `json:"value"`
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

func FilterSpotsByFeature(feats []FeatureFilter) (*[]Spot, error) {
	var spots []Spot

	query := Database.Joins("JOIN spots_to_feats ON spots_to_feats.id = spots.id")

	// Append Where conditions for each feature
	for _, feat := range feats {
		if feat.Value != nil {
			// If feature is with a value
			query = query.Where("spots_to_feats.feat_id = ? AND spots_to_feats.value = ?", feat.ID, *feat.Value)
		} else {
			// If feature is without the value
			query = query.Where("spots_to_feats.feat_id = ? AND spots_to_feats.value IS NULL", feat.ID)
		}
	}

	err := query.
		Group("spot_id").
		Having("COUNT(*) = ?", len(feats)).
		Preload("Features.Feature").
		Find(&spots).Error

	if err != nil {
		return nil, err
	}

	return &spots, nil

}

/*
type Feat struct {
	ID       uint   `json:id`
	Value    *int8  `json:"value"`
}
*/
