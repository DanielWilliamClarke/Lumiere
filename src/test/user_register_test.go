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
	"dwc.com/lumiere/mocks"
	"dwc.com/lumiere/user"
)

func RunRegisterTest(mocks registerMocks, body []byte) *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest("POST", "http://localhost:5000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Setup Test router
	app := fiber.New()
	app.Post("/", user.UserRegisterRoute{
		DataAccess: mocks.mockClient,
		Generator:  mocks.mockGenerator,
	}.Post)

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

type registerMocks struct {
	mockClient    *mocks.MockIMongoClient
	mockGenerator *mocks.MockICodeGenerator
	ctrl          *gomock.Controller
}

type registerMockConfig struct {
	t               *testing.T
	findError       error
	insertError     error
	totalFindCalls  int
	generatorErrors []error
}

func ConfigureRegisterMocks(config registerMockConfig) registerMocks {

	ctrl := gomock.NewController(config.t)

	mockClient := mocks.NewMockIMongoClient(ctrl)
	mockGenerator := mocks.NewMockICodeGenerator(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(config.totalFindCalls).
		Return(config.findError)

	// Set up generator expectations
	generatorWillError := false
	if config.findError != nil {
		generatorCalls := make([]*gomock.Call, 0)
		for _, err := range config.generatorErrors {

			code := "test-code"
			if err != nil {
				generatorWillError = true
				code = ""
			}

			generatorCalls = append(generatorCalls, mockGenerator.EXPECT().
				Generate(gomock.AssignableToTypeOf(0)).
				Times(1).
				Return(code, err))

			// Sorry - if the an error appears then we wont be calling generate any more times
			if err != nil {
				break
			}
		}
		gomock.InOrder(generatorCalls...)
	} else {
		mockGenerator.EXPECT().
			Generate(gomock.AssignableToTypeOf(0)).
			Times(0)
	}

	// If generator will error we will not be inserting
	totalInsertCalls := 1
	if config.findError == nil || generatorWillError {
		totalInsertCalls = 0
	}

	mockClient.EXPECT().
		InsertOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf([]byte{})).
		Times(totalInsertCalls).
		Return(nil, config.insertError)

	return registerMocks{mockClient, mockGenerator, ctrl}
}

func Test_UserCanBeRegistered(t *testing.T) {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test
	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       errors.New("No data found"),
		insertError:     nil,
		totalFindCalls:  1,
		generatorErrors: []error{nil, nil},
	})
	defer mocks.ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunRegisterTest(mocks, body)
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status 200")
	}
}

func Test_UserInsertFails(t *testing.T) {

	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       errors.New("No data found"),
		insertError:     errors.New("Fake insert failure"),
		totalFindCalls:  1,
		generatorErrors: []error{nil, nil},
	})

	defer mocks.ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunRegisterTest(mocks, body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}

func Test_UserFindFails(t *testing.T) {

	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       nil,
		insertError:     nil,
		totalFindCalls:  1,
		generatorErrors: []error{nil, nil},
	})
	defer mocks.ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunRegisterTest(mocks, body)
	if resp.StatusCode != http.StatusNoContent {
		t.Error("Expected status 204")
	}
}

func Test_UserCantBeRegisteredIfUserInputIsInvalid(t *testing.T) {

	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       nil,
		insertError:     nil,
		totalFindCalls:  0,
		generatorErrors: []error{nil, nil},
	})
	defer mocks.ctrl.Finish()

	resp := RunRegisterTest(mocks, []byte{})
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Expected status 400")
	}
}

func Test_UserIDGenerationFails(t *testing.T) {

	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       errors.New("No data found"),
		insertError:     nil,
		totalFindCalls:  1,
		generatorErrors: []error{nil, errors.New("user id generation failure")},
	})
	defer mocks.ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunRegisterTest(mocks, body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}

func Test_UserCodeGenerationFails(t *testing.T) {

	mocks := ConfigureRegisterMocks(registerMockConfig{
		t:               t,
		findError:       errors.New("No data found"),
		insertError:     nil,
		totalFindCalls:  1,
		generatorErrors: []error{errors.New("user code generation failure"), nil},
	})
	defer mocks.ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunRegisterTest(mocks, body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}
