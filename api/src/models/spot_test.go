package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("POSTGRES_URL", "postgresql://localhost:5432/hot-chip_test")
	OpenDatabaseConnection()
	AutoMigrateModels()
	os.Exit(m.Run())
}

func cleanupSpots() {
	Database.Exec("TRUNCATE TABLE spots;")
}

func TestSpot_Save(t *testing.T) {
	cleanupSpots()
	spot := Spot{
		Name:    "Test Cafe",
		Address: "1 Test Street",
	}

	savedSpot, err := spot.Save()

	assert.Nil(t, err)
	assert.Equal(t, "Test Cafe", savedSpot.Name)
	assert.NotZero(t, savedSpot.ID)
}

func TestFetchAllSpots(t *testing.T) {
	cleanupSpots()
	spot1 := Spot{Name: "Cafe One", Address: "1 Test Street"}
	spot2 := Spot{Name: "Cafe Two", Address: "2 Test Street"}
	spot1.Save()
	spot2.Save()

	spots, err := FetchAllSpots()

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(*spots), 2)
}

func TestFindSpot(t *testing.T) {
	cleanupSpots()
	spot := Spot{Name: "Find This Cafe", Address: "3 Test Street"}
	savedSpot, _ := spot.Save()

	foundSpot, err := FindSpot(fmt.Sprintf("%d", savedSpot.ID))

	assert.Nil(t, err)
	assert.Equal(t, "Find This Cafe", foundSpot.Name)
}
