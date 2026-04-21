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
