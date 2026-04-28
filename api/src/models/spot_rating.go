package models

import "gorm.io/gorm"

type SpotRating struct {
	gorm.Model
	UserId uint  `json:"user_id"`
	SpotId uint  `json:"spot_id"`
	Rating int8  `json:"rating"`
}

// AddRating adds a new rating or updates an existing one if user has already rated this spot
func AddRating(userId uint, spotId uint, rating int8) error {
	spotRating := SpotRating{
		UserId: userId,
		SpotId: spotId,
		Rating: rating,
	}

	// upsert - find existing rating for this user+spot
	// update it or create a new one if it doesn't exist 
	err := Database.
		Where(SpotRating{UserId: userId, SpotId: spotId}).
		Assign(SpotRating{Rating: rating}).
		FirstOrCreate(&spotRating).Error

	return err
}

// returns average rating for a single spot
// returns nil if no ratings exist yet instead of 0
func GetAvgRatingForSpotById(spotId uint) (*float64, error) {
	var avg *float64
	err := Database.Model(&SpotRating{}).
		Where("spot_id = ?", spotId).
		Select("AVG(rating)").
		Scan(&avg).Error
	return avg, err
}

// returns map of SpotId and average rating for all spots 
// spots with no ratings wont appear in the map 
func GetAvgRatingsForSpots() (map[uint]*float64, error) {
	var results []struct {
		SpotId uint
		Avg    *float64
	}
	err := Database.Model(&SpotRating{}).
		Select("spot_id, AVG(rating) as avg").
		Group("spot_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// convert slice of results into a map for faster lookup by spot id 
	avgMap := make(map[uint]*float64)
	for _, r := range results {
		avgMap[r.SpotId] = r.Avg
	}
	return avgMap, nil
}