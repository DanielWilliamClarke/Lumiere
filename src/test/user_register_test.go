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

func RunTest(mockClient lMongo.IMongoClient, body []byte) *http.Response {
	// Create request
	req := httptest.NewRequest("POST", "http://localhost:5000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Setup Test router
	app := fiber.New()
	app.Post("/", user.UserRegisterRoute{DataAccess: mockClient}.Post)

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

func ConfigureMocks(t *testing.T, findError error, insertError error) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

	ctrl := gomock.NewController(t)
	mockClient := mock_mongo.NewMockIMongoClient(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(1).
		Return(findError)

	totalInsertCalls := 1
	if findError == nil {
		totalInsertCalls = 0
	}

	mockClient.EXPECT().
		InsertOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf([]byte{})).
		Times(totalInsertCalls).
		Return(nil, insertError)

	return mockClient, ctrl
}

func Test_AccountCanBeRegistered(t *testing.T) {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	mockClient, ctrl := ConfigureMocks(t, errors.New("No data found"), nil)
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTest(mockClient, body)
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}

func Test_AccountInsertFails(t *testing.T) {

	mockClient, ctrl := ConfigureMocks(t, errors.New("No data found"), errors.New("Fake insert failure"))
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTest(mockClient, body)
	if resp.StatusCode != 500 {
		t.Error("Expected status 500")
	}
}

func Test_AccountFindFails(t *testing.T) {

	mockClient, ctrl := ConfigureMocks(t, nil, nil)
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTest(mockClient, body)
	if resp.StatusCode != 500 {
		t.Error("Expected status 500")
	}
}
