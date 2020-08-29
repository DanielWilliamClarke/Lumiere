package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	"dwc.com/lumiere/mocks"
	mock_mongo "dwc.com/lumiere/mocks"
	lMongo "dwc.com/lumiere/mongo"
	"dwc.com/lumiere/user"
)

func RunAuthTest(mockClient lMongo.IMongoClient, auth string) *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest("GET", "http://localhost:5000", nil)
	req.Header.Set("Authorization", auth)
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

func ConfigureAuthMocks(t *testing.T, findError error, totalCalls int) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockIMongoClient(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(totalCalls).
		Return(findError)

	return mockClient, ctrl
}

func Test_UserCanBeAuthorized(t *testing.T) {

	mockClient, ctrl := ConfigureAuthMocks(t, nil, 1)
	defer ctrl.Finish()

	resp := RunAuthTest(mockClient, "test")
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}

func Test_UserCanBeUnauthorized(t *testing.T) {

	mockClient, ctrl := ConfigureAuthMocks(t, errors.New("User not found"), 1)
	defer ctrl.Finish()

	resp := RunAuthTest(mockClient, "test")
	if resp.StatusCode != 403 {
		t.Error("Expected status 403")
	}
}

func Test_UserIsUnauthorizedWithEmptyUserCode(t *testing.T) {

	mockClient, ctrl := ConfigureAuthMocks(t, nil, 0)
	defer ctrl.Finish()

	resp := RunAuthTest(mockClient, "")
	if resp.StatusCode != 403 {
		t.Error("Expected status 403")
	}
}
