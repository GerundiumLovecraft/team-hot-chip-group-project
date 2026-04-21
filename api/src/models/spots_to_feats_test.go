package models

import (
	"testing"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSTFModelSuiteEnv struct {
	suite.Suite
	db *gorm.DB
}

// Tests are run before they start
func (suite *TestSTFModelSuiteEnv) SetupSuite() {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	suite.db = Database
}

func (suite *TestSTFModelSuiteEnv) SetupTest() {
	spot1 := Spot{
		Name:    "spot1",
		Address: "Spot Street",
	}
	spot2 := Spot{
		Name:    "spot2",
		Address: "Sponge Street",
	}

	spot1.Save()
	spot2.Save()

	feature1 := Feature{
		FeatName: "great place",
	}
	feature2 := Feature{
		FeatName: "amazing people",
	}

	feature1.SaveNewFeature()
	feature2.SaveNewFeature()
}

// Running after each test
func (suite *TestSTFModelSuiteEnv) TearDownTest() {
	suite.db.Raw(`TRUNCATE TABLE spots;`)
	suite.db.Raw(`TRUNCATE TABLE features;`)
	suite.db.Raw(`TRUNCATE TABLE spots_to_features;`)
}

func TestSTFModelSuite(t *testing.T) {
	suite.Run(t, new(TestSTFModelSuiteEnv))
}

func (suite *TestSTFModelSuiteEnv) TestSaveNewRelation() {
	
}
