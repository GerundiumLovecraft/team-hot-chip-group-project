package models

type SpotsToFeats struct {
	SpotId  uint    `gorm:"primaryKey;autoIncrement:false"`
	FeatId  uint    `gorm:"primaryKey;autoIncrement:false" json:"feat_id"`
	Value   *int8   `gorm:"check:value IS NULL OR (value>=1 AND value<=3)"`
	Spot    Spot    `gorm:"foreignKey:SpotId"`
	Feature Feature `gorm:"foreignKey:FeatId"`
}

func (spotToFeat *SpotsToFeats) SaveNewRelation() (*SpotsToFeats, error) {
	err := Database.Create(spotToFeat).Error

	if err != nil {
		return &SpotsToFeats{}, err
	}

	return spotToFeat, nil
}

func (spotToFeat *SpotsToFeats) UpdateRelation(value int8) error {
	prevValue := spotToFeat.Value
	spotToFeat.Value = &value

	err := Database.Save(spotToFeat).Error

	if err != nil {
		spotToFeat.Value = prevValue
		return err
	}

	return nil
}

func DeleteSpotToFeat(spotId uint, featId uint) error {
	return Database.Delete(&SpotsToFeats{}, "spot_id = ? AND feat_id = ?", spotId, featId).Error
}
