package models

type Feature struct {
	ID       uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	FeatName string `gorm:"type:varchar(255);not null" json:"feat_name"`
}

func FetchAllFeatures() (*[]Feature, error) {
	var features []Feature

	err := Database.Find(&features).Error

	if err != nil {
		return nil, err
	}

	return &features, nil
}

func (f *Feature) SaveNewFeature() (*Feature, error) {
	err := Database.Create(f).Error

	if err != nil {
		return &Feature{}, err
	}

	return f, nil
}

func SeedFeatures() { 
	// check if features already exist in the database
	var count int64
	Database.Model(&Feature{}).Count(&count)

	// if features already exist, don't add duplicates
	if count > 0 {
		return
	}

	// features to seed
	features := []Feature{
		{FeatName: "wifi"},
		{FeatName: "toilets"},
		{FeatName: "power_sockets"},
		{FeatName: "open_late"},
		{FeatName: "noise_level"},
		{FeatName: "price"},
	}

	// save each feature to the database
	for _, feature := range features {
		f := feature 
		f.SaveNewFeature()
	}
}
