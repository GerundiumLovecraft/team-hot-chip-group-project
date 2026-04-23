// Credentials to log in during development:
// bambi@sipstack.com / password123
// thumper@sipstack.com / password123
// Note: seed runs automatically on startup if no spots exist
// To reseed from scratch run: TRUNCATE TABLE spots_to_feats, spots, users CASCADE; then restart server

package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func SeedDemoData() {
	// check if spots already exist — skip if already seeded
	var count int64
	Database.Model(&Spot{}).Count(&count)
	if count > 0 {
		fmt.Println("Demo data already exists, skipping seed")
		return
	}

	// hash password for all demo users
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// create demo user 1
	user1 := User{
		Username:       "bambi",
		Email:          "bambi@sipstack.com",
		HashedPassword: string(hashedPassword),
	}
	savedUser1, err := user1.Save()
	if err != nil {
		fmt.Println("Error creating bambi:", err)
		return
	}

	// create demo user 2
	user2 := User{
		Username:       "thumper",
		Email:          "thumper@sipstack.com",
		HashedPassword: string(hashedPassword),
	}
	savedUser2, err := user2.Save()
	if err != nil {
		fmt.Println("Error creating thumper:", err)
		return
	}

	fmt.Println("Demo users created")

	// create demo spots spread across users
	spots := []Spot{
		{
			UserId:      savedUser1.ID,
			Name:        "Ducie Street Warehouse",
			Address:     "Ducie Street, Manchester, M1 2TP",
			Description: "A Grade II listed Victorian warehouse just 100 yards from Piccadilly Station. By day it fills with remote workers settling into dark wood booths with laptops. Great coffee counter, lounge seating, and an outdoor terrace.",
			OpenFrom:    "07:30",
			OpenTo:      "18:00",
		},
		{
			UserId:      savedUser1.ID,
			Name:        "Haunt",
			Address:     "58 Peter Street, Manchester, M2 3NQ",
			Description: "A fiercely independent specialty coffee shop on Peter Street with a Mediterranean soul. Monochrome chequered floors, Bauhaus art prints and free wifi make it a popular haunt for remote workers. Opens at 7:30am on weekdays.",
			OpenFrom:    "07:30",
			OpenTo:      "22:00",
		},
		{
			UserId:      savedUser2.ID,
			Name:        "WatchHouse Marylebone",
			Address:     "32-34 New Cavendish Street, London, W1G 8UE",
			Description: "A bright two-floor specialty coffee shop in the heart of Marylebone Village. Known for exceptional single origin coffee and a full brunch menu. Quieter downstairs section is ideal for laptop work. WiFi available with a 90 minute limit.",
			OpenFrom:    "07:30",
			OpenTo:      "18:00",
		},
		{
			UserId:      savedUser2.ID,
			Name:        "Books N' Cup Cafe",
			Address:     "23 Home Street, Edinburgh, EH3 9JR",
			Description: "A charming Edinburgh sanctuary for book lovers and coffee enthusiasts. Rustic decor, cozy leather seating, shelves of vintage books and artisanal ceramic cups. Some tables are designated laptop-free. Open late — perfect for evening working sessions.",
			OpenFrom:    "09:00",
			OpenTo:      "22:00",
		},
	}

	for _, spot := range spots {
		s := spot
		savedSpot, err := s.Save()
		if err != nil {
			fmt.Println("Error creating spot:", err)
			return
		}

		// link wifi feature to each spot
		wifiFeature := SpotsToFeats{
			SpotId: savedSpot.ID,
			FeatId: 1,
		}
		wifiFeature.SaveNewRelation()
	}

	fmt.Println("Demo spots created successfully")
}