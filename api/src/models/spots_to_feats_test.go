package models

import (
	"testing"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSTFModelSuiteEnv struct {
	suite.Suite
	db      *gorm.DB
	spotId1 uint
	featId1 uint
}

// Tests are run before they start
func (suite *TestSTFModelSuiteEnv) SetupSuite() {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	suite.db = Database
}

func (suite *TestSTFModelSuiteEnv) SetupTest() {

	feature1 := Feature{

		FeatName: "great place",
	}

	savedFeature1, _ := feature1.SaveNewFeature()
	suite.featId1 = savedFeature1.ID

	spot1 := Spot{
		Name:    "spot1",
		Address: "Spot Street",
	}

	savedSpot1, _ := spot1.Save()
	suite.spotId1 = savedSpot1.ID

}

// Running after each test
func (suite *TestSTFModelSuiteEnv) TearDownTest() {
	suite.db.Exec(`TRUNCATE spots, features, spots_to_feats CASCADE;`)
}

func TestSTFModelSuite(t *testing.T) {
	suite.Run(t, new(TestSTFModelSuiteEnv))
}

func (suite *TestSTFModelSuiteEnv) TestSaveNewRelation() {
	newRelation := SpotsToFeats{
		SpotId: suite.spotId1,
		FeatId: suite.featId1,
	}

	savedRel, err := newRelation.SaveNewRelation()

	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Equal(suite.T(), newRelation.SpotId, savedRel.SpotId, "Saved relation should be equal to new relation")
	assert.Equal(suite.T(), newRelation.FeatId, savedRel.FeatId, "Saved relation should be equal to new relation")
}

func (suite *TestSTFModelSuiteEnv) TestUpdateRelation() {
	var relValue int8 = 1

	newRelation := SpotsToFeats{
		SpotId: suite.spotId1,
		FeatId: suite.featId1,
		Value:  &relValue,
	}
	// happy path
	newRelation.SaveNewRelation()
	err1 := newRelation.UpdateRelation(3)

	assert.Nil(suite.T(), err1, "Spot Id should be updated")
	assert.Equal(suite.T(), int8(3), *newRelation.Value, "Value field should be updated")

	// unhappy path
	relValue = 4

	err2 := newRelation.UpdateRelation(relValue)

	assert.NotNil(suite.T(), err2, "Error should be thrown if updated value is invalid")
	assert.Equal(suite.T(), int8(3), *newRelation.Value, "Value field should be reverted to previous state")
}

func (suite *TestSTFModelSuiteEnv) TestDeleteRelation() {
	var relValue int8 = 1

	var relation1 SpotsToFeats = SpotsToFeats{
		SpotId: suite.spotId1,
		FeatId: suite.featId1,
		Value:  &relValue,
	}

	relation1.SaveNewRelation()

	err := DeleteSpotToFeat(suite.spotId1, suite.featId1)

	assert.Nil(suite.T(), err, "Error should be nil")
}
