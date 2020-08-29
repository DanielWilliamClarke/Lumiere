package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	mock_mongo "dwc.com/lumiere/mocks"
	lMongo "dwc.com/lumiere/mongo"
	"dwc.com/lumiere/user"
)

func RunAuthTest(mockClient lMongo.IMongoClient, body []byte) *http.Response {
	// Create request
	req := httptest.NewRequest("GET", "http://localhost:5000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Setup Test router
	app := fiber.New()
	app.Use(user.UserAuthMiddleware{DataAccess: mockClient}.Auth).
		Get("/", func(c *fiber.Ctx) {
			c.Status(200).Send("Authorized")
		})

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

func ConfigureAuthMocks(t *testing.T, findError error) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

	ctrl := gomock.NewController(t)
	mockClient := mock_mongo.NewMockIMongoClient(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(1).
		Return(findError)

	return mockClient, ctrl
}

func Test_UserCanBeAuthorized(t *testing.T) {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	mockClient, ctrl := ConfigureAuthMocks(t, nil)
	defer ctrl.Finish()

	body, err := json.Marshal(user.UserCodeBody{
		UserCode: "test",
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunAuthTest(mockClient, body)
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}

func Test_UserCanBeUnauthorized(t *testing.T) {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	mockClient, ctrl := ConfigureAuthMocks(t, errors.New("User not found"))
	defer ctrl.Finish()

	body, err := json.Marshal(user.UserCodeBody{
		UserCode: "test",
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunAuthTest(mockClient, body)
	if resp.StatusCode != 403 {
		t.Error("Expected status 403")
	}
}
