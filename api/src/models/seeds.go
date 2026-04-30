// Credentials for all demo users: username@sipstack.com / password123
// To reseed from scratch: TRUNCATE TABLE spots_to_feats, spots, users CASCADE; then restart server

package models

import (
	"fmt"
	"net/url"
	"golang.org/x/crypto/bcrypt"
)

func buildEmbedURL(name string, address string) string {
	query := url.QueryEscape(name + " " + address)
	return fmt.Sprintf("https://www.google.com/maps?q=%s&z=17&output=embed", query)
}

func SeedDemoData() {
	var count int64
	Database.Model(&Spot{}).Count(&count)
	if count > 0 {
		fmt.Println("Demo data already exists, skipping seed")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// ── Users ──────────────────────────────────────────────────────────────────
	users := []User{
		{
			Username:       "wednaguirand",
			Email:          "wedna@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "",
		},
		{
			Username:       "shakiraw",
			Email:          "shakira@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "",
		},
		{
			Username:       "jackgarner",
			Email:          "jack@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "",
		},
		{
			Username:       "vladk",
			Email:          "vladislav@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "",
		},
		{
			Username:       "bobbelcher",
			Email:          "bobby@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "https://i.pinimg.com/736x/5c/5b/75/5c5b754ab2e9fe8884448ce9743a7041.jpg",
		},
		{
			Username:       "dwight",
			Email:          "dwight@sipstack.com",
			HashedPassword: string(hashedPassword),
			Avatar:         "https://i.pinimg.com/1200x/cb/89/dd/cb89dde2bc82ee72ed23d3ba0d88a6cb.jpg",
		},
	}

	savedUsers := make([]*User, len(users))
	for i, u := range users {
		user := u
		saved, err := user.Save()
		if err != nil {
			fmt.Printf("Error creating user %s: %v\n", user.Username, err)
			return
		}
		savedUsers[i] = saved
	}

	fmt.Println("Demo users created")

	quiet      := int8(1)
	moderate   := int8(2)
	onePrice   := int8(1)
	twoPrice   := int8(2)
	threePrice := int8(3)

	type spotSeed struct {
		spot     Spot
		features func(spotId uint)
	}

	spotSeeds := []spotSeed{

		// 0 — Manchester ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Foundation Coffee House",
				Address:     "48-50 Whitworth Street, Manchester, M1 6LS",
				Description: "One of Manchester's most loved independent coffee shops, Foundation is spacious, bright and genuinely built for remote workers. Excellent flat whites powered by local roasters, a full breakfast and lunch menu, and healthy snack options mean there's no need to leave until you're done. The open layout means you never feel cramped and the high window seats are perfect for a change of scenery. Fast WiFi, plenty of sockets and a relaxed attitude towards laptops.",
				OpenFrom: "07:30", OpenTo: "20:00",
				Image: "https://i.pinimg.com/736x/a3/e0/be/a3e0be2b0c9ff90086c1b6d6529da400.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 1 — London ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[1].ID, Name: "Prufrock Coffee",
				Address:     "23-25 Leather Lane, London, EC1N 7TE",
				Description: "One of London's most celebrated specialty coffee destinations, Prufrock on Leather Lane is also an excellent place to work. Small tables, natural light and seriously good coffee make it ideal for focused sessions. It sits on Leather Lane Market which runs every weekday lunchtime with dozens of international street food stalls — meaning lunch is sorted without going far. The brunch menu inside is equally excellent.",
				OpenFrom: "08:00", OpenTo: "19:00",
				Image: "https://i.pinimg.com/736x/3c/01/df/3c01dffff34f15a71542f62352277af8.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 2 — Edinburgh ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[2].ID, Name: "Cairngorm Coffee",
				Address:     "1 Melville Place, Edinburgh, EH3 7PR",
				Description: "One of Edinburgh's best-known work-friendly cafes, Cairngorm on Melville Street is bright, spacious and never makes you feel guilty about staying all afternoon. The coffee is exceptional and some tables even have built-in iPads for quick tasks. The WiFi is fast and free and the communal tables are ideal for solo workers. A genuinely brilliant place to spend a full working day in Edinburgh.",
				OpenFrom: "07:00", OpenTo: "19:00",
				Image: "https://i.pinimg.com/1200x/52/c0/21/52c02107f85e4508272da0ad01ad33aa.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 3 — Liverpool
		{
			Spot{
				UserId: savedUsers[3].ID, Name: "Bold Street Coffee",
				Address:     "89 Bold Street, Liverpool, L1 4HF",
				Description: "Bold Street Coffee is a Liverpool institution. Even on a quiet weekday morning there's an energy to this place that makes it inspiring to work in. The menu is excellent — the hash browns in particular have developed a cult following — and the coffee is reliably great. It's small so arrive early or on a weekday afternoon for the best chance of a seat. The atmosphere more than makes up for the size.",
				OpenFrom: "08:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/14/8e/54/148e54584048921f8394ec6fc5e4928e.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 4 — Manchester
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Takk",
				Address:     "6 Tariff Street, Manchester, M1 2FF",
				Description: "A slice of Scandinavia in Manchester's Northern Quarter. Takk is beloved by creatives and freelancers for its relaxed atmosphere, excellent coffee and Nordic-inspired interiors. The back section has a dedicated row of two-person tables with sockets beneath each one — perfect for a long session. Their sourdough and pastries are excellent and the flat whites are consistently some of the best in the city. Dogs welcome too.",
				OpenFrom: "08:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/97/bf/5e/97bf5ef9f89b96252b72833d5320fe6c.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 5 — London
		{
			Spot{
				UserId: savedUsers[1].ID, Name: "Monmouth Coffee Company",
				Address:     "27 Monmouth Street, London, WC2H 9EU",
				Description: "A London institution that has been at the forefront of specialty coffee since 1978, Monmouth in Covent Garden is a favourite of students and freelancers who like working in a constant energising hum of conversation. The menu is short and focused — no frills, just perfectly executed espresso, pour-over and filter coffee. A little pricier than average but the quality is unmatched. Best visited on a quieter weekday morning.",
				OpenFrom: "08:00", OpenTo: "18:30",
				Image: "https://i.pinimg.com/736x/6e/9d/c3/6e9dc30c75f5b621db66f3ad111cbcc0.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &threePrice})
			},
		},

		// 6 — Edinburgh
		{
			Spot{
				UserId: savedUsers[2].ID, Name: "Brew Lab",
				Address:     "6-8 South College Street, Edinburgh, EH8 9AA",
				Description: "A staple of Edinburgh's specialty coffee scene, Brew Lab has a rustic aesthetic of distressed wood and exposed brick that feels genuinely lived-in. Popular with students and professionals alike. Their espresso-based drinks are carefully prepared and the oat milk flat whites are a highlight. It's lively during term time which creates a productive buzz — best for working sessions of a few hours. The vegan donut options are not to be missed.",
				OpenFrom: "08:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/a0/a1/b6/a0a1b6c3c9e0600ffd54322e48efd5d1.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 7 — Liverpool
		{
			Spot{
				UserId: savedUsers[3].ID, Name: "Press Bros Coffee",
				Address:     "14 Harrington Street, Liverpool, L2 9QA",
				Description: "A quiet gem in Liverpool city centre, Press Bros is the kind of place where you can set up for the morning and genuinely get things done. The coffee is carefully sourced and the service is warm and unhurried. The food menu is simple but good — pastries, toasties and a daily soup. Sockets are available throughout and the WiFi is dependable. A reliable and underrated option for remote workers in Liverpool.",
				OpenFrom: "07:30", OpenTo: "16:00",
				Image: "https://i.pinimg.com/736x/9e/87/11/9e871181cf47a92b012db175cafb9077.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &onePrice})
			},
		},

		// 8 — Manchester
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Haunt",
				Address:     "58 Peter Street, Manchester, M2 3NQ",
				Description: "A fiercely independent specialty coffee shop and wine bar on Peter Street. Monochrome chequered floors, Bauhaus prints and a Mediterranean soul make it a genuinely atmospheric place to work. Popular with remote workers who don't want to feel like they're in a co-working space. Open until 11pm most evenings — one of the few Manchester spots where you can shift from a productive afternoon into an evening glass of natural wine without moving.",
				OpenFrom: "07:30", OpenTo: "23:00",
				Image: "https://i.pinimg.com/736x/4f/8b/32/4f8b327bafa6efc06d7fa1f2813efe9c.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 9 — London
		{
			Spot{
				UserId: savedUsers[1].ID, Name: "Origin Coffee Shoreditch",
				Address:     "65 Charlotte Road, London, EC2A 3PE",
				Description: "Origin's Shoreditch outpost is a strong choice for remote workers who need space. Large communal tables make it easy to spread out and the coffee sourcing is among the best in East London. The food covers pastries and light bites throughout the day. A great option for both solo sessions and collaborative working — the atmosphere is creative without being too loud.",
				OpenFrom: "07:30", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/ad/f6/e4/adf6e4233d42f51879d46c20e25b2270.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 10 — Edinburgh
		{
			Spot{
				UserId: savedUsers[2].ID, Name: "Black Sheep Coffee",
				Address:     "25 North Bridge, Edinburgh, EH1 1SB",
				Description: "Conveniently located on North Bridge close to Waverley station, Black Sheep Coffee is one of Edinburgh's more dependable work spots. Open until 8pm most evenings — a solid option if you need to work later than most cafes allow. The WiFi is consistently good, the team are laptop-friendly and the coffee is solid. Not the most atmospheric spot in Edinburgh but reliably comfortable and brilliantly located for getting things done.",
				OpenFrom: "07:00", OpenTo: "20:00",
				Image: "https://i.pinimg.com/1200x/e4/57/09/e45709de168bb6f73aaf658b1a058332.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 11 — Manchester
		{
			Spot{
				UserId: savedUsers[4].ID, Name: "Foley's Cafe",
				Address:     "3 Byrom Street, Manchester, M3 4PP",
				Description: "Set in Manchester's Spinningfields district, Foley's is purpose-built for remote workers and digital nomads. The seating is intelligently arranged — large booths, communal tables and a mezzanine level with individual desks overlooking the cafe. Fast WiFi, reliable sockets and a strong phone signal make it one of the most practical spots in the city. The food menu covers pastries, sandwiches, salad boxes and smoothies — no need to leave until you're ready.",
				OpenFrom: "07:30", OpenTo: "17:00",
				Image: "https://i.pinimg.com/1200x/6f/61/07/6f610755f065c00bdd7b82392b6d68b1.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 12 — London ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[1].ID, Name: "Timberyard Seven Dials",
				Address:     "7 Upper St Martin's Lane, London, WC2H 9DL",
				Description: "A consistent favourite among London's remote working crowd, Timberyard in Seven Dials has built its reputation on being genuinely laptop-friendly. Good sockets, reliable WiFi and a calm atmosphere make it easy to settle in for hours. The coffee is sourced carefully and the food menu covers everything from morning pastries to a solid lunch. One of the more polished work-friendly cafes in central London.",
				OpenFrom: "07:30", OpenTo: "20:00",
				Image: "https://i.pinimg.com/736x/10/74/3c/10743c9ad7bdab07a6f5a2bcf1003f75.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 13 — Edinburgh
		{
			Spot{
				UserId: savedUsers[2].ID, Name: "Fortitude Coffee",
				Address:     "3C York Place, Edinburgh, EH1 3EB",
				Description: "A beloved Edinburgh institution near York Place, Fortitude is small but mighty. The single-origin coffees are exceptional and among the most affordable in the city — a flat white for under £3 is rare these days. The baristas are knowledgeable and genuinely passionate about their craft. Best for solo sessions of an hour or two rather than full-day camps, but a wonderful spot to start your working morning before moving on.",
				OpenFrom: "07:30", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/bc/19/32/bc193206fde6b0fae455ab2f8291965f.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &onePrice})
			},
		},

		// 14 — Liverpool
		{
			Spot{
				UserId: savedUsers[3].ID, Name: "Ryde",
				Address:     "27 Hope Street, Liverpool, L1 9BQ",
				Description: "A relaxed and creative coffee shop on Hope Street, sitting between Liverpool's two cathedrals in one of the city's most atmospheric streets. The interiors are warm and eclectic, the coffee is excellent and the kitchen turns out genuinely good food throughout the day. A welcoming spot for remote workers with enough space to spread out and a vibe that manages to be both productive and enjoyable.",
				OpenFrom: "08:00", OpenTo: "17:30",
				Image: "https://i.pinimg.com/736x/cd/4a/2f/cd4a2f5733fca44a5d5867a127ff260c.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 15 — Manchester
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Pot Kettle Black",
				Address:     "Barton Arcade, 51-63 Deansgate, Manchester, M3 2BH",
				Description: "Set inside the stunning Victorian Barton Arcade in Deansgate, Pot Kettle Black brings Australian and New Zealand coffee culture to Manchester. The setting is beautiful — high glass ceilings, ornate ironwork — and the coffee is exceptional. Best for focused morning sessions before the lunch crowd arrives. The avocado toast is worth ordering and the baristas are always happy to talk about their beans.",
				OpenFrom: "08:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/69/49/bb/6949bb15fac0408d0551df6f5a8fa1c5.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 16 — London
		{
			Spot{
				UserId: savedUsers[1].ID, Name: "Hackney Coffee Company",
				Address:     "503 Hackney Road, London, E2 9ED",
				Description: "With its exposed yellow stock bricks and industrial interiors, Hackney Coffee Company is a characterful spot that works surprisingly well as a remote office. Popular with East London's creative community, there's a lovely courtyard garden for warmer months. Come the evening, soft lighting and candles transform the space — one of the few spots that transitions beautifully from a working afternoon into an evening out. The coffee is excellent and the team are warm.",
				OpenFrom: "08:00", OpenTo: "22:00",
				Image: "https://i.pinimg.com/1200x/33/ad/5b/33ad5b0c0a81b5cb5a039712bd4dd5b0.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 17 — Manchester
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Ezra & Gil",
				Address:     "2 Stevenson Square, Manchester, M1 1DN",
				Description: "Winner of the 2025 I Love Manchester Most Loved Coffee Shop Award, Ezra & Gil on Stevenson Square lives up to the hype. The relaxed neighbourhood feel makes it easy to settle in for a long session. Their brunch menu is creative and genuinely delicious — think French toast and seasonal specials — and the flat whites are powered by ManCoCo beans. Fast WiFi, enough sockets and staff who won't rush you out.",
				OpenFrom: "08:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/0e/bd/4d/0ebd4d96497be9b4c354c1bf05932ef9.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &moderate})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 18 — Edinburgh ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[4].ID, Name: "Black Medicine Coffee Co",
				Address:     "2 Nicolson Street, Edinburgh, EH8 9DH",
				Description: "Black Medicine on Nicolson Street is one of Edinburgh's most popular spots for deep focus working. The large basement area is ideal for spreading out with minimal distractions and a quiet atmosphere that encourages concentration. The coffee is decent and the vegan menu is a highlight — including a notoriously good brownie. The team are relaxed about how long you stay. Best in the morning before the student crowd arrives.",
				OpenFrom: "08:00", OpenTo: "20:00",
				Image: "https://i.pinimg.com/736x/40/20/44/402044ec401931fab63020f87a296f13.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 19 — Manchester
		{
			Spot{
				UserId: savedUsers[0].ID, Name: "Fig & Sparrow",
				Address:     "20 Oldham Street, Manchester, M1 1JN",
				Description: "An icon of the Northern Quarter, Fig & Sparrow has a cosy, almost cottage-like atmosphere that makes it one of Manchester's nicest spots to settle in with a laptop. The menu is seasonal and interesting, the coffee is excellent and the booths give you a sense of your own corner. Quieter than most Northern Quarter spots and a reliable choice for a focused morning session. The homemade cakes are hard to resist.",
				OpenFrom: "09:00", OpenTo: "17:00",
				Image: "https://i.pinimg.com/736x/bb/0b/c8/bb0bc80158755a04cefcedebcc1c18f9.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &onePrice})
			},
		},

		// 20 — Manchester ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[4].ID, Name: "Chapter One Books",
				Address:     "Chatsworth House, 23 Lever Street, Manchester, M1 1BY",
				Description: "An independent bookshop and café in Manchester's Northern Quarter that started from a family's love of great books, good coffee and homemade cakes. The atmosphere is warm and unhurried making it one of the easier places in the city to lose yourself in a few hours of focused work. Open until 9pm most evenings — one of the few spots where you can comfortably work into the evening with a great slice of cake for company.",
				OpenFrom: "09:00", OpenTo: "21:00",
				Image: "https://i.pinimg.com/736x/60/9b/7b/609b7b8b5f769361a41a3a2681c4144a.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},

		// 21 — London
		{
			Spot{
				UserId: savedUsers[5].ID, Name: "WatchHouse Marylebone",
				Address:     "32-34 New Cavendish Street, London, W1G 8UE",
				Description: "A bright two-floor specialty coffee shop in the heart of Marylebone Village, WatchHouse is known for exceptional single origin coffee and a beautifully curated brunch menu. The quieter downstairs section is ideal for laptop work. A premium experience — the coffee is among the best in London and the food is thoughtfully put together. One of the most polished and reliable work-friendly cafes in the capital.",
				OpenFrom: "07:30", OpenTo: "16:30",
				Image: "https://i.pinimg.com/736x/a9/b2/c1/a9b2c14f00e1020cdb78d26a7e077e43.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &threePrice})
			},
		},

		// 22 — London ✅ wifi, toilets, quiet, ££, open_late
		{
			Spot{
				UserId: savedUsers[5].ID, Name: "FWD Coffee",
				Address:     "54 Farringdon Road, London, EC1R 3BL",
				Description: "Formerly known as Powerhouse Coffee, FWD on Farringdon Road is a specialty coffee shop with a warm team and a strong reputation among London's remote working crowd. The coffee is sourced carefully and executed with real attention. A comfortable and focused atmosphere that suits long working sessions, with enough sockets and reliable WiFi to keep you going all day. A solid and underrated option in Clerkenwell.",
				OpenFrom: "07:30", OpenTo: "19:00",
				Image: "https://i.pinimg.com/736x/d4/9e/3f/d49e3fe3a8c13a43e157c01c4dc2afa5.jpg",
			},
			func(id uint) {
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 1})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 2})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 3})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 4})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 5, Value: &quiet})
				Database.Create(&SpotsToFeats{SpotId: id, FeatId: 6, Value: &twoPrice})
			},
		},
	}

	// Save all spots, attach features and build map URL
	savedSpots := make([]*Spot, len(spotSeeds))
	for i, s := range spotSeeds {
		spot := s.spot
		spot.LocationURL = buildEmbedURL(spot.Name, spot.Address)
		saved, err := spot.Save()
		if err != nil {
			fmt.Printf("Error creating spot %s: %v\n", spot.Name, err)
			return
		}
		s.features(saved.ID)
		savedSpots[i] = saved
	}

	fmt.Println("Demo spots created successfully")

	// ── Ratings ────────────────────────────────────────────────────────────────
	type ratingEntry struct {
		userId  uint
		spotIdx int
		rating  int8
	}

	ratings := []ratingEntry{
		// Five-star (15)
		{savedUsers[1].ID, 0, 5},  
		{savedUsers[2].ID, 4, 5},  
		{savedUsers[3].ID, 8, 5},  
		{savedUsers[4].ID, 15, 5}, 
		{savedUsers[5].ID, 17, 5}, 
		{savedUsers[0].ID, 1, 5},  
		{savedUsers[2].ID, 5, 5}, 
		{savedUsers[3].ID, 9, 5},  
		{savedUsers[4].ID, 12, 5}, 
		{savedUsers[5].ID, 16, 5},
		{savedUsers[0].ID, 2, 5}, 
		{savedUsers[1].ID, 6, 5}, 
		{savedUsers[0].ID, 3, 5}, 
		{savedUsers[1].ID, 7, 5}, 
		{savedUsers[2].ID, 11, 5},
		// Four-star (5)
		{savedUsers[2].ID, 19, 4}, 
		{savedUsers[1].ID, 10, 4}, 
		{savedUsers[4].ID, 13, 4}, 
		{savedUsers[0].ID, 14, 4}, 
		{savedUsers[0].ID, 21, 4}, 
		// Three-star (3)
		{savedUsers[2].ID, 18, 3}, 
		{savedUsers[5].ID, 3, 3}, 
		{savedUsers[3].ID, 19, 3},
		// Extra
		{savedUsers[5].ID, 20, 5}, 
		{savedUsers[3].ID, 22, 5},
	}

	for _, r := range ratings {
		err := AddRating(r.userId, savedSpots[r.spotIdx].ID, r.rating)
		if err != nil {
			fmt.Printf("Error adding rating: %v\n", err)
		}
	}

	fmt.Println("Demo ratings added successfully")
}