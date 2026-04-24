// Credentials to log in:
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
			Image:       "https://images.pexels.com/photos/36038130/pexels-photo-36038130.jpeg",
		},
		{
			UserId:      savedUser1.ID,
			Name:        "Haunt",
			Address:     "58 Peter Street, Manchester, M2 3NQ",
			Description: "A fiercely independent specialty coffee shop on Peter Street with a Mediterranean soul. Monochrome chequered floors, Bauhaus art prints and free wifi make it a popular haunt for remote workers. Opens at 7:30am on weekdays.",
			OpenFrom:    "07:30",
			OpenTo:      "22:00",
			Image:       "https://images.pexels.com/photos/33456309/pexels-photo-33456309.jpeg",
		},
		{
			UserId:      savedUser2.ID,
			Name:        "WatchHouse Marylebone",
			Address:     "32-34 New Cavendish Street, London, W1G 8UE",
			Description: "A bright two-floor specialty coffee shop in the heart of Marylebone Village. Known for exceptional single origin coffee and a full brunch menu. Quieter downstairs section is ideal for laptop work. WiFi available with a 90 minute limit.",
			OpenFrom:    "07:30",
			OpenTo:      "16:30",
			Image:       "https://images.pexels.com/photos/33505093/pexels-photo-33505093.jpeg",
		},
		{
			UserId:      savedUser2.ID,
			Name:        "Books N' Cup Cafe",
			Address:     "23 Home Street, Edinburgh, EH3 9JR",
			Description: "A charming Edinburgh sanctuary for book lovers and coffee enthusiasts. Rustic decor, cozy leather seating, shelves of vintage books and artisanal ceramic cups. Some tables are designated laptop-free. Open late — perfect for evening working sessions.",
			OpenFrom:    "09:00",
			OpenTo:      "22:00",
			Image:       "https://images.pexels.com/photos/8820015/pexels-photo-8820015.jpeg",
		},
	}

	for i, spot := range spots {
		s := spot
		savedSpot, err := s.Save()
		if err != nil {
			fmt.Println("Error creating spot:", err)
			return
		}

		switch i {
		case 0: // Ducie Street Warehouse
			moderate := int8(2)
			twoPrice := int8(2)
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 1})                     // wifi
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 2})                     // toilets
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 3})                     // power_sockets
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 5, Value: &moderate})   // noise_level: moderate
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 6, Value: &twoPrice})   // price: ££

		case 1: // Haunt
			moderate := int8(2)
			twoPrice := int8(2)
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 1})                     // wifi
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 2})                     // toilets
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 4})                     // open_late
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 5, Value: &moderate})   // noise_level: moderate
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 6, Value: &twoPrice})   // price: ££

		case 2: // WatchHouse Marylebone
			quiet := int8(1)
			threePrice := int8(3)
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 1})                     // wifi
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 2})                     // toilets
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 3})                     // power_sockets
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 5, Value: &quiet})      // noise_level: quiet
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 6, Value: &threePrice}) // price: £££

		case 3: // Books N' Cup Cafe
			quiet := int8(1)
			onePrice := int8(1)
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 1})                     // wifi
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 4})                     // open_late
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 5, Value: &quiet})      // noise_level: quiet
			Database.Create(&SpotsToFeats{SpotId: savedSpot.ID, FeatId: 6, Value: &onePrice})   // price: £
		}
	}

	fmt.Println("Demo spots created successfully")
}