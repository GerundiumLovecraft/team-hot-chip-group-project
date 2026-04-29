package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gorm.io/gorm"
)

type TestSuiteEnv struct {
	suite.Suite
	db    *gorm.DB
	token string
	app   *gin.Engine
	res   *httptest.ResponseRecorder
	ids   struct {
		spotId1 uint
		spotId2 uint
		featId1 uint
	}
}

func (suite *TestSuiteEnv) SignupAndLogin(email, username, password string) string {
	suite.res = httptest.NewRecorder()
	signupJSON := []byte(fmt.Sprintf(
		`{"email":"%s", "password":"%s", "username":"%s"}`,
		email, password, username,
	))
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJSON))
	suite.app.ServeHTTP(suite.res, signupReq)
	assert.Equal(suite.T(), 201, suite.res.Code)

	suite.res = httptest.NewRecorder()
	loginJSON := []byte(fmt.Sprintf(
		`{"usernameOrEmail":"%s", "password":"%s"}`,
		email, password,
	))
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJSON))
	suite.app.ServeHTTP(suite.res, loginReq)

	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	return loginResponse["token"]
}

func (suite *TestSuiteEnv) CreateSpot(token, name string) {
	suite.res = httptest.NewRecorder()
	spotJSON := []byte(fmt.Sprintf(
		`{"user_id": 1, "name": "%s", "address": "main_test.go", "description": "test desc", "open_from": "11:00", "open_to": "19:00", "features": [{"feat_id": %d}]}`,
		name, suite.ids.featId1,
	))
	postSpotRequest, _ := http.NewRequest("POST", "/spots", bytes.NewBuffer(spotJSON))
	postSpotRequest.Header.Set("Authorization", "Bearer "+token)
	suite.app.ServeHTTP(suite.res, postSpotRequest)
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	env.LoadEnv(".test.env")
	models.OpenDatabaseConnection()
	models.AutoMigrateModels()
	models.SeedFeatures()
	suite.db = models.Database
	suite.app = setupApp()
	suite.token, _ = auth.GenerateToken("test-user")

}

func (suite *TestSuiteEnv) SetupTest() {
	suite.res = httptest.NewRecorder()

	spot1 := models.Spot{
		Name:    "spot1",
		Address: "main_test.go",
	}
	spot2 := models.Spot{
		Name:    "spot2",
		Address: "main_test.go",
	}
	savedSpot1, _ := spot1.Save()
	savedSpot2, _ := spot2.Save()

	feat1 := models.Feature{
		FeatName: "feat1",
	}

	savedFeat1, _ := feat1.SaveNewFeature()

	suite.ids.spotId1 = savedSpot1.ID
	suite.ids.spotId2 = savedSpot2.ID
	suite.ids.featId1 = savedFeat1.ID
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	suite.db.Exec("TRUNCATE TABLE users CASCADE;")
	suite.db.Exec("TRUNCATE TABLE posts CASCADE;")
	suite.db.Exec("TRUNCATE TABLE spots CASCADE;")
	suite.db.Exec("TRUNCATE TABLE features CASCADE;")
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv) Test_PostUsers_IncorrectJSON() {
	app, token := suite.app, suite.token

	var jsonStr = []byte(`{"message":"Test Post"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

// =============================================
// SIGNUP INTEGRATION TESTS

func (suite *TestSuiteEnv) Test_SignupUser_CorrectCredentials() {
	// valid signup should return 201
	app, token := suite.app, suite.token

	var jsonStr = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 201, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_SignupUser_MissingCredentials() {
	// missing email and password should return 400
	app, token := suite.app, suite.token

	var jsonStr = []byte(`{"email":"", "password":"", "username":"testuser"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_SignupUser_DuplicateEmail() {
	// duplicate email should return 409
	app, token := suite.app, suite.token

	// first signup
	var jsonStr = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	// second signup with same email
	var jsonStr2 = []byte(`{"email":"test@example.com", "password":"password123", "username":"differentuser"}`)
	suite.res = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr2))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 409, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_SignupUser_DuplicateUsername() {
	// duplicate username should return 409
	app, token := suite.app, suite.token

	// first signup
	var jsonStr = []byte(`{"email":"test@example.com", "password":"password123", "username":"sameuser"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	// second signup with same username
	var jsonStr2 = []byte(`{"email":"different@example.com", "password":"password123", "username":"sameuser"}`)
	suite.res = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr2))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 409, suite.res.Code)
}

// =============================================
// LOGIN INTEGRATION TESTS

func (suite *TestSuiteEnv) Test_LoginUser_CorrectCredentials() {
	// create a user first so we have someone to log in as
	app := suite.app
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	// login with same credentials
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	// expect 201 response back
	assert.Equal(suite.T(), 201, suite.res.Code)

	// is response data correct
	var response map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &response)
	assert.NotEmpty(suite.T(), response["token"])
}

func (suite *TestSuiteEnv) Test_LoginUser_WrongPassword() {
	app := suite.app
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password1234"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_LoginUser_EmailNotFound() {
	app := suite.app

	var loginJson = []byte(`{"usernameOrEmail":"notregistered@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

// =============================================
// SPOTS INTEGRATION TESTS

// Test GET /spots
func (suite *TestSuiteEnv) Test_GetSpots_ReturnsListOfSpots() {
	// Send GET request to /spots
	app := suite.app
	getSpotsRequest, _ := http.NewRequest("GET", "/spots", nil)
	app.ServeHTTP(suite.res, getSpotsRequest)

	// Retrieve information from response
	var response struct {
		Spots []struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"spots"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &response)

	// Assert the results
	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), suite.ids.spotId1, response.Spots[0].ID)
	assert.Equal(suite.T(), suite.ids.spotId2, response.Spots[1].ID)
	assert.Equal(suite.T(), 2, len(response.Spots))
}

// Test GET /spots/:id
func (suite *TestSuiteEnv) Test_GetSpotById_ReturnsSpot() {
	// Send GET request to
	app := suite.app
	getSpotRequest, _ := http.NewRequest("GET", fmt.Sprintf("/spots/%v", suite.ids.spotId1), nil)
	app.ServeHTTP(suite.res, getSpotRequest)

	// Retrieve information from the response
	var response struct {
		Spot struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"Spot"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &response)

	// Assert the results
	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), suite.ids.spotId1, response.Spot.ID)
	assert.Equal(suite.T(), "spot1", response.Spot.Name)
	assert.Equal(suite.T(), "main_test.go", response.Spot.Address)
}

func (suite *TestSuiteEnv) Test_GetSpotById_ReturnsError_IncorrectId() {
	// Send GET request to
	app := suite.app
	getSpotRequest, _ := http.NewRequest("GET", "/spots/9999999", nil)
	app.ServeHTTP(suite.res, getSpotRequest)

	assert.Equal(suite.T(), 404, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_CreateSpot_AuthError_WithoutToken() {
	// Send a POST request to /spots without token
	app := suite.app

	spotJson := []byte(`{"user_id": "1", "name": "spot3", "address": "main_test.go", "description": "test desc", "open_from": "11:00", "open_to": "19:00", "features": { "feat_name": "test_feature"}}`)
	postSpotRequest, _ := http.NewRequest("POST", "/spots", bytes.NewBuffer(spotJson))
	postSpotRequest.Header.Set("Authorization", "Bearer undefined")
	app.ServeHTTP(suite.res, postSpotRequest)

	// Assert the error code
	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_CreateSpot_BadReq_IncompleteBody() {
	//signup
	app := suite.app
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)
	//login
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	var response map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &response)
	// Send POST request to /spots with incomplete body
	suite.res = httptest.NewRecorder()
	spotJson := []byte(`{"user_id": "1", "name": "spot3", "description": "test desc", "open_from": "11:00", "open_to": "19:00", "features": { "feat_name": "test_feature"}}`)
	postSpotRequest, _ := http.NewRequest("POST", "/spots", bytes.NewBuffer(spotJson))
	postSpotRequest.Header.Set("Authorization", "Bearer "+response["token"])
	app.ServeHTTP(suite.res, postSpotRequest)

	// Assert the error code
	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_CreateSpot_OK() {
	//signup
	app := suite.app
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)
	//login
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)
	// Send POST request to /spots with incomplete body
	suite.res = httptest.NewRecorder()
	spotJson := []byte(fmt.Sprintf(`{"user_id": 1, "name": "spot3", "address": "main_test.go", "description": "test desc", "open_from": "11:00", "open_to": "19:00", "features": [{"feat_id": %d}]}`, suite.ids.featId1))
	postSpotRequest, _ := http.NewRequest("POST", "/spots", bytes.NewBuffer(spotJson))
	postSpotRequest.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, postSpotRequest)

	//retrieve information from the response
	var postResponse struct {
		SpotID uint   `json:"spotID"`
		Token  string `json:"token"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &postResponse)

	// Assert the result
	assert.Equal(suite.T(), 201, suite.res.Code)
	assert.NotZero(suite.T(), postResponse.SpotID)
}

// =============================================
// PROFILE INTEGRATION TESTS
// Tests GET /users/:id

func (suite *TestSuiteEnv) Test_GetUserByID_OwnProfile() {
	app := suite.app

	// signup
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	// login to get a real token with the users actual ID
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	// get token from the login response
	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// get the users ID from the database
	user, _ := models.FindUserByUsernameOrEmail("test@example.com")
	userId := fmt.Sprintf("%d", user.ID)

	// request own profile
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId, nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 200, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetUserByID_Not_OwnProfile() {
	app := suite.app

	// person 1 sign up
	var signupJson1 = []byte(`{"email":"person1@example.com", "password":"password123", "username":"person1"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson1))
	app.ServeHTTP(suite.res, signupReq)

	// person 2 sign up
	var signupJson2 = []byte(`{"email":"person2@example.com", "password":"password123", "username":"person2"}`)
	signupReq, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson2))
	app.ServeHTTP(suite.res, signupReq)

	// person 1 login and get token
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"person1@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// get person 2 user ID from the database
	user, _ := models.FindUserByUsernameOrEmail("person2@example.com")
	userId := fmt.Sprintf("%d", user.ID)

	// request person 2 profile
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+userId, nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetUserByID_NotFound() {
	app := suite.app

	// signup
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	// login to get a real token with the users actual ID
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	// get token from the login response
	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// request a fake user ID that doesn't exist
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/99999", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 404, suite.res.Code)
}

// Tests GET /users/search-by-username/:username

func (suite *TestSuiteEnv) Test_GetUserByUsername_OwnProfile() {
	app := suite.app

	// signup
	var signupJson = []byte(`{"email":"test@example.com", "password":"password123", "username":"testuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	// login to get a real token
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"test@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	// get token from the login response
	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// request own profile
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search-by-username/testuser", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 200, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetUserByUsername_Not_OwnProfile() {
	app := suite.app

	// person 1 sign up
	var signupJson1 = []byte(`{"email":"person1@example.com", "password":"password123", "username":"person1"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson1))
	app.ServeHTTP(suite.res, signupReq)

	// person 2 sign up
	var signupJson2 = []byte(`{"email":"person2@example.com", "password":"password123", "username":"person2"}`)
	signupReq, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson2))
	app.ServeHTTP(suite.res, signupReq)

	// person 1 login and get token
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"person1@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// request person 2 profile
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search-by-username/person2", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetUserByUsername_NotFound() {
	app := suite.app

	// signup
	var signupJson = []byte(`{"email":"realuser@example.com", "password":"password123", "username":"realuser"}`)
	signupReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(signupJson))
	app.ServeHTTP(suite.res, signupReq)

	// login to get a real token
	suite.res = httptest.NewRecorder()
	var loginJson = []byte(`{"usernameOrEmail":"realuser@example.com", "password":"password123"}`)
	loginReq, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(loginJson))
	app.ServeHTTP(suite.res, loginReq)

	// get token from the login response
	var loginResponse map[string]string
	json.Unmarshal(suite.res.Body.Bytes(), &loginResponse)

	// request a fake username that doesn't exist
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/search-by-username/idontexist", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse["token"])
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 404, suite.res.Code)
}

// =============================================
// FEATURES INTEGRATION TESTS
// GET /features

func (suite *TestSuiteEnv) Test_GetAllFeatures_ReturnsAllFeatures() {
	app := suite.app

	getFeatsRequest, _ := http.NewRequest("GET", "/features", nil)
	app.ServeHTTP(suite.res, getFeatsRequest)

	// Retrieve information from the response
	var response struct {
		Features []struct {
			ID   uint   `json:"feat_id"`
			Name string `json:"feat_name"`
		}
	}

	json.Unmarshal(suite.res.Body.Bytes(), &response)

	// Assert the results
	assert.Equal(suite.T(), 200, suite.res.Code, "Reponse code should be 200")
	assert.Equal(suite.T(), 1, len(response.Features), "Reponse length should be 7")

}

// =============================================
// LEADERBOARDS INTEGRATION TESTS
// GET /leaderboards

func (suite *TestSuiteEnv) Test_GetLeaderboardOrderedBySpots() {
	app := suite.app

	type leaderboardResponse struct {
		Leaderboard []struct {
			UserID       uint   `json:"user_id"`
			Username     string `json:"username"`
			SpotsCreated int    `json:"spots_created"`
		} `json:"leaderboard"`
	}

	tokenA := suite.SignupAndLogin("a@example.com", "usera", "password123")
	suite.CreateSpot(tokenA, "a-spot-1")
	suite.CreateSpot(tokenA, "a-spot-2")

	tokenB := suite.SignupAndLogin("b@example.com", "userb", "password123")
	suite.CreateSpot(tokenB, "b-spot-1")

	tokenC := suite.SignupAndLogin("c@example.com", "userc", "password123")
	suite.CreateSpot(tokenC, "c-spot-1")
	suite.CreateSpot(tokenC, "c-spot-2")
	suite.CreateSpot(tokenC, "c-spot-3")

	//Leaderboard get request
	suite.res = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/leaderboard", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	app.ServeHTTP(suite.res, req)

	var response leaderboardResponse
	err := json.Unmarshal(suite.res.Body.Bytes(), &response)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Len(suite.T(), response.Leaderboard, 3)

	//Assert in descending order
	assert.Equal(suite.T(), "userc", response.Leaderboard[0].Username)
	assert.Equal(suite.T(), 3, response.Leaderboard[0].SpotsCreated)
	assert.NotZero(suite.T(), response.Leaderboard[0].UserID)

	assert.Equal(suite.T(), "usera", response.Leaderboard[1].Username)
	assert.Equal(suite.T(), 2, response.Leaderboard[1].SpotsCreated)
	assert.NotZero(suite.T(), response.Leaderboard[1].UserID)

	assert.Equal(suite.T(), "userb", response.Leaderboard[2].Username)
	assert.Equal(suite.T(), 1, response.Leaderboard[2].SpotsCreated)
	assert.NotZero(suite.T(), response.Leaderboard[2].UserID)

	for i := 0; i < len(response.Leaderboard)-1; i++ {
		assert.GreaterOrEqual(
			suite.T(),
			response.Leaderboard[i].SpotsCreated,
			response.Leaderboard[i+1].SpotsCreated,
		)
	}
}

/*
GET /profile

	case 1: Valid token → 200
	case 2: Invalid/No token → 401
*/
func (suite *TestSuiteEnv) Test_GetProfile_ValidToken() {
	app := suite.app

	// Create user A
	tokenA := suite.SignupAndLogin("a@example.com", "usera", "password123")

	suite.res = httptest.NewRecorder()

	// Construct a GET /profile request
	getProfileRequest, _ := http.NewRequest("GET", "/profile", nil)
	getProfileRequest.Header.Set("Authorization", "Bearer "+tokenA)
	app.ServeHTTP(suite.res, getProfileRequest)

	// Retrieve information from the response
	var getProfileResponse struct {
		User struct {
			ID        uint      `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"createdAt"`
			Avatar    string    `json:"avatar"`
		} `json:"user"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &getProfileResponse)
	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.NotNil(suite.T(), getProfileResponse)
	assert.Equal(suite.T(), "a@example.com", getProfileResponse.User.Email)
}

func (suite *TestSuiteEnv) Test_GetProfile_InvalidToken() {
	app := suite.app

	// Create invalid token
	tokenA := "invalid-token"

	suite.res = httptest.NewRecorder()

	// Construct a GET /profile request
	getProfileRequest, _ := http.NewRequest("GET", "/profile", nil)
	getProfileRequest.Header.Set("Authorization", "Bearer "+tokenA)
	app.ServeHTTP(suite.res, getProfileRequest)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

/*
GET /profile/spots
	case 1: Valid Token + Spots > 200 + slice of spots
	case 2: Valid Token + No Spots > 200 + empty slice
	case 3: Invalid/No Token > 401
*/

func (suite *TestSuiteEnv) Test_GetSpotsByUser_WithSpots() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")
	suite.CreateSpot(token, "spot-1")
	suite.CreateSpot(token, "spot-2")

	app := suite.app
	suite.res = httptest.NewRecorder()

	getSpotsByUserRequest, _ := http.NewRequest("GET", "/profile/spots", nil)
	getSpotsByUserRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, getSpotsByUserRequest)

	var getSpotsByUserResponse struct {
		Spots []struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"spots"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &getSpotsByUserResponse)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), 2, len(getSpotsByUserResponse.Spots))
	assert.Equal(suite.T(), "spot-1", getSpotsByUserResponse.Spots[0].Name)

}

func (suite *TestSuiteEnv) Test_GetSpotsByUser_WithNoSpots() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	getSpotsByUserRequest, _ := http.NewRequest("GET", "/profile/spots", nil)
	getSpotsByUserRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, getSpotsByUserRequest)

	var getSpotsByUserResponse struct {
		Spots []struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"spots"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &getSpotsByUserResponse)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), 0, len(getSpotsByUserResponse.Spots))
}

func (suite *TestSuiteEnv) Test_GetSpotsByUser_InvalidToken() {
	token := "invalid-token"

	app := suite.app
	suite.res = httptest.NewRecorder()

	getSpotsByUserRequest, _ := http.NewRequest("GET", "/profile/spots", nil)
	getSpotsByUserRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, getSpotsByUserRequest)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

/*
POST /spots/filter

	case 1: Valid feature filter that matches spots → 200 + matching spots returned
	case 2: Valid feature filter with no matches → 200 + empty array + message "Spots not found"
	case 3: Empty array body [] → 400 "Filter parameters are empty"
	case 4: Malformed JSON body → 400
*/
func (suite *TestSuiteEnv) Test_GetSpotsByFeature_Match() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")
	suite.CreateSpot(token, "spot-1")
	suite.CreateSpot(token, "spot-2")

	app := suite.app
	suite.res = httptest.NewRecorder()

	var featureJson = []byte(fmt.Sprintf(`[{"feat_id": %v}]`, suite.ids.featId1))
	postSpotsByFeatureRequest, _ := http.NewRequest("POST", "/spots/filter", bytes.NewBuffer(featureJson))
	app.ServeHTTP(suite.res, postSpotsByFeatureRequest)

	var postSpotsByFeatureResponse struct {
		Spots []struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"spots"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &postSpotsByFeatureResponse)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), 2, len(postSpotsByFeatureResponse.Spots))
	assert.Equal(suite.T(), "spot-1", postSpotsByFeatureResponse.Spots[0].Name)
}

func (suite *TestSuiteEnv) Test_GetSpotsByFeature_NoMatch() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")
	suite.CreateSpot(token, "spot-1")
	suite.CreateSpot(token, "spot-2")

	app := suite.app
	suite.res = httptest.NewRecorder()

	var featureJson = []byte(`[{"feat_id": 9876}]`)
	postSpotsByFeatureRequest, _ := http.NewRequest("POST", "/spots/filter", bytes.NewBuffer(featureJson))
	app.ServeHTTP(suite.res, postSpotsByFeatureRequest)

	var postSpotsByFeatureResponse struct {
		Spots []struct {
			ID      uint   `json:"_id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"spots"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &postSpotsByFeatureResponse)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), 0, len(postSpotsByFeatureResponse.Spots))
}

func (suite *TestSuiteEnv) Test_GetSpotsByFeature_EmptyFeatArray() {
	app := suite.app
	suite.res = httptest.NewRecorder()

	var featureJson = []byte(`[]`)
	postSpotsByFeatureRequest, _ := http.NewRequest("POST", "/spots/filter", bytes.NewBuffer(featureJson))
	app.ServeHTTP(suite.res, postSpotsByFeatureRequest)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetSpotsByFeature_MalformedJson() {
	app := suite.app
	suite.res = httptest.NewRecorder()

	var featureJson = []byte(`not a valid json at all`)
	postSpotsByFeatureRequest, _ := http.NewRequest("POST", "/spots/filter", bytes.NewBuffer(featureJson))
	app.ServeHTTP(suite.res, postSpotsByFeatureRequest)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

/*
POST /spots/:id/rate
	case 1: Valid token + valid spot ID + rating between 1-5 → 201
	case 2: No token → 401
	case 3: Rating below 1 or above 5 → 400
	case 4: Missing rating field in body → 400
	case 5: Invalid spot ID (non-numeric) → 400
*/

func (suite *TestSuiteEnv) Test_AddRating() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	var ratingJson = []byte(`{"rating": 4}`)
	postNewRatingRequest, _ := http.NewRequest("POST", fmt.Sprintf(`/spots/%v/rate`, suite.ids.spotId1), bytes.NewBuffer(ratingJson))
	postNewRatingRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, postNewRatingRequest)

	assert.Equal(suite.T(), 201, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_AddRating_NoToken() {
	app := suite.app
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	var ratingJson = []byte(`{"rating": 4}`)
	postNewRatingRequest, _ := http.NewRequest("POST", fmt.Sprintf(`/spots/%v/rate`, suite.ids.spotId1), bytes.NewBuffer(ratingJson))
	app.ServeHTTP(suite.res, postNewRatingRequest)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_AddRating_InvalidRating() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app

	// Send a request with a rating lower than 1
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	var lowRatingJson = []byte(`{"rating": 0}`)
	postNewRatingRequestLow, _ := http.NewRequest("POST", fmt.Sprintf(`/spots/%v/rate`, suite.ids.spotId1), bytes.NewBuffer(lowRatingJson))
	postNewRatingRequestLow.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, postNewRatingRequestLow)

	assert.Equal(suite.T(), 400, suite.res.Code)

	// Send a request with a rating higher than 5
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	var highRatingJson = []byte(`{"rating": 6}`)
	postNewRatingRequestHigh, _ := http.NewRequest("POST", fmt.Sprintf(`/spots/%v/rate`, suite.ids.spotId1), bytes.NewBuffer(highRatingJson))
	postNewRatingRequestHigh.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, postNewRatingRequestHigh)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_AddRating_NoRating() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	postNewRatingRequestLow, _ := http.NewRequest("POST", fmt.Sprintf(`/spots/%v/rate`, suite.ids.spotId1), nil)
	postNewRatingRequestLow.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, postNewRatingRequestLow)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_AddRating_InvalidSpotId() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	// Construct a POST request to add a review
	var ratingJson = []byte(`{"rating": 4}`)
	postNewRatingRequestLow, _ := http.NewRequest("POST", `/spots/invalid-id/rate`, bytes.NewBuffer(ratingJson))
	postNewRatingRequestLow.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, postNewRatingRequestLow)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

/*
PATCH /users/avatar
	case 1: Valid token + valid avatar URL → 200 + updated user returned
	case 2: No token → 401
	case 3: Missing avatar field in body → 400
*/

func (suite *TestSuiteEnv) Test_UpdateUserAvatar() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	var avatarJson = []byte(`{"avatar": "https://marketplace.canva.com/EAFdIOc_EM0/3/0/1600w/canva-beige-and-pink-illustrative-cute-girl-avatar-pUk78cukytY.jpg"}`)
	patchUserAvatarRequest, _ := http.NewRequest("PATCH", "/users/avatar", bytes.NewBuffer(avatarJson))
	patchUserAvatarRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, patchUserAvatarRequest)

	var patchUpdateAvatarRequest struct {
		User struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			CreatedAt string `json:"createdAt"`
			Avatar    string `json:"avatar"`
		} `json:"user"`
	}
	json.Unmarshal(suite.res.Body.Bytes(), &patchUpdateAvatarRequest)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), "usera", patchUpdateAvatarRequest.User.Username)
	assert.NotNil(suite.T(), patchUpdateAvatarRequest.User.Avatar)
}

func (suite *TestSuiteEnv) Test_UpdateUserAvatar_NoToken() {
	app := suite.app
	suite.res = httptest.NewRecorder()

	var avatarJson = []byte(`{"avatar": "https://marketplace.canva.com/EAFdIOc_EM0/3/0/1600w/canva-beige-and-pink-illustrative-cute-girl-avatar-pUk78cukytY.jpg"}`)
	patchUserAvatarRequest, _ := http.NewRequest("PATCH", "/users/avatar", bytes.NewBuffer(avatarJson))
	app.ServeHTTP(suite.res, patchUserAvatarRequest)

	assert.Equal(suite.T(), 401, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_UpdateUserAvatar_MissingJsonBody() {
	token := suite.SignupAndLogin("a@example.com", "usera", "password123")

	app := suite.app
	suite.res = httptest.NewRecorder()

	patchUserAvatarRequest, _ := http.NewRequest("PATCH", "/users/avatar", nil)
	patchUserAvatarRequest.Header.Set("Authorization", "Bearer "+token)
	app.ServeHTTP(suite.res, patchUserAvatarRequest)

	assert.Equal(suite.T(), 400, suite.res.Code)
}
