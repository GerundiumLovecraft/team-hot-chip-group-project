package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
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

func (suite *TestSuiteEnv) Test_PostUsers_CorrectJSON() {
	app, token := suite.app, suite.token

	res := httptest.NewRecorder()
	var jsonStr = []byte(`{"email":"test@email.com", "password": "testpassword"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(res, req)

	assert.Equal(suite.T(), 201, res.Code)
}

func (suite *TestSuiteEnv) Test_PostUsers_IncorrectJSON() {
	app, token := suite.app, suite.token

	var jsonStr = []byte(`{"message":"Test Post"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	assert.Equal(suite.T(), 400, suite.res.Code)
}

func (suite *TestSuiteEnv) Test_GetPosts() {
	app, token := suite.app, suite.token

	newPost := models.Post{
		Message: "Test Post",
	}
	newPost.Save()

	req, _ := http.NewRequest("GET", "/posts", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	app.ServeHTTP(suite.res, req)

	responseData, _ := io.ReadAll(suite.res.Body)
	var jsonPosts struct {
		Posts []controllers.JSONPost
	}

	_ = json.Unmarshal(responseData, &jsonPosts)

	assert.Equal(suite.T(), 200, suite.res.Code)
	assert.Equal(suite.T(), "Test Post", jsonPosts.Posts[0].Message)
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
