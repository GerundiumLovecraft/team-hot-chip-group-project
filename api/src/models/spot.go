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
	ID    uint  `json:"feat_id"`
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

	query := Database.Table("spots_to_feats").
		Select("spot_id").
		Group("spot_id").
		Having("COUNT(*) = ?", len(feats))


	// Append Where conditions for each feature
	for _, feat := range feats {
		if feat.Value != nil {
			// If feature is with a value
			query = query.Or("feat_id = ? AND value = ?", feat.ID, *feat.Value)
		} else {
			// If feature is without the value
			query = query.Or("feat_id = ? AND value IS NULL", feat.ID)
		}
	}

	err := Database.Preload("Features.Feature").
		Where("id IN (?)", query).
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
	Avatar       string `json:"avatar"`
}

func FetchLeaderboard() ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := Database.Table("spots").
	    Select("spots.user_id, users.username, users.avatar, COUNT(spots.id) as spots_created, MIN(spots.created_at) as first_spot_at").
        Joins("JOIN users ON users.id = spots.user_id").
        Group("spots.user_id, users.username, users.avatar").
        Order("spots_created DESC, first_spot_at ASC").
        Scan(&entries).Error

	if err != nil {
		return []LeaderboardEntry{}, err
	}

	return entries, nil
}