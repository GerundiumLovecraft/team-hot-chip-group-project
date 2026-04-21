package models

import (
	"testing"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestFeatureModelSuiteEnv struct {
	suite.Suite
	db *gorm.DB
}

// Tests are run before they start
func (suite *TestFeatureModelSuiteEnv) SetupSuite() {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	suite.db = Database
}

// Running after each test
func (suite *TestFeatureModelSuiteEnv) TearDownTest() {
	suite.db.Raw("TRUNCATE TABLE features;")
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestFeatureModelSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(TestFeatureModelSuiteEnv))
}

func (suite *TestFeatureModelSuiteEnv) TestSaveNewFeature() {
	feature := Feature{
		FeatName: "test",
	}

	savedFeature, err := feature.SaveNewFeature()

	assert.Nil(suite.T(), err, "Error should be nil")
	assert.NotEmpty(suite.T(), savedFeature, "Saved feature should not be empty")
}

func (suite *TestFeatureModelSuiteEnv) TestFetchAllFeatures() {
	// setup the database
	feature1 := Feature{
		FeatName: "test-1",
	}
	feature2 := Feature{
		FeatName: "test-2",
	}
	feature1.SaveNewFeature()
	feature2.SaveNewFeature()

	features, err := FetchAllFeatures()

	assert.Nil(suite.T(), err, "Error should be nil")
	assert.Len(suite.T(), *features, 2, "Number of features should be 2")

}
