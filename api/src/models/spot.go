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
	Image       string         `json:"image"`
	OpenFrom    string         `json:"open_from"`
	OpenTo      string         `json:"open_to"`
	Features    []SpotsToFeats `json:"features" gorm:"foreignKey:SpotId;references:ID;constraint:OnDelete:CASCADE"`
	LocationURL	string			`json:"location_url"`
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

	var featuresIds []uint
	for _, feat := range feats {
		featuresIds = append(featuresIds, feat.ID)
	}

	query := Database.Table("spots").Joins("JOIN spots_to_feats ON spots_to_feats.spot_id = spots.id").Joins("JOIN features ON features.id = spots_to_feats.feat_id").Where("features.id IN ?", featuresIds)

	err := query.
		Group("spots.id").
		Having("COUNT(DISTINCT spots_to_feats.feat_id) = ?", len(featuresIds)).
		Preload("Features.Feature").
		Find(&spots).Error

	if err != nil {
		return nil, err
	}

	return &spots, nil

}

func FetchSpotsByUserId(userId string) (*[]Spot, error) {
    var spots []Spot
    err := Database.Preload("Features.Feature").Where("user_id = ?", userId).Find(&spots).Error
    if err != nil {
        return &[]Spot{}, err
    }
    return &spots, nil
}

/*
type Feat struct {
	ID       uint   `json:id`
	Value    *int8  `json:"value"`
}
*/

type LeaderboardEntry struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	SpotsCreated int    `json:"spots_created" gorm:"column:spots_created"`
}

func FetchLeaderboard() ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := Database.Table("spots").
	    Select("spots.user_id, users.username, COUNT(spots.id) as spots_created, MIN(spots.created_at) as first_spot_at").
        Joins("JOIN users ON users.id = spots.user_id").
        Group("spots.user_id, users.username").
        Order("spots_created DESC, first_spot_at ASC").
        Scan(&entries).Error

	if err != nil {
		return []LeaderboardEntry{}, err
	}

	return entries, nil
}