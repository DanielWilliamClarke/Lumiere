package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang/mock/gomock"
	"github.com/ulule/deepcopier"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account"
	"dwc.com/lumiere/account/model"
	"dwc.com/lumiere/mocks"
	mock_mongo "dwc.com/lumiere/mocks"
	lMongo "dwc.com/lumiere/mongo"
	"dwc.com/lumiere/user"
)

func RunTransferTest(mockClient lMongo.IMongoClient, authedAccount interface{}, body []byte) *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Setup Test router
	app := fiber.New()
	app.
		Use(func(c *fiber.Ctx) {
			c.Locals("AuthedAccount", authedAccount)
			c.Next()
		}).
		Put("/", account.AccountTransferRoute{DataAccess: mockClient}.PutTransfer)

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

type transferMockConfig struct {
	t              *testing.T
	findError      error
	findData       interface{}
	totalFindCalls int
	unaryTransfer  bool
	updateErrors   []error
}

func ConfigureTransferMocks(config transferMockConfig) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

	ctrl := gomock.NewController(config.t)
	mockClient := mocks.NewMockIMongoClient(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(config.totalFindCalls).
		DoAndReturn(func(ctx context.Context, filter interface{}, data interface{}) error {
			if config.findData != nil {
				deepcopier.Copy(config.findData).To(data)
			}
			return config.findError
		})

	if config.totalFindCalls > 0 && config.findError == nil && !config.unaryTransfer {

		mockClient.EXPECT().
			StartTransaction(
				gomock.AssignableToTypeOf(context.Background()),
				gomock.AssignableToTypeOf(func() error { return nil })).
			Times(1).
			DoAndReturn(func(ctx context.Context, callback func() error) error {
				return callback()
			})

		updateCalls := make([]*gomock.Call, 0)
		for _, err := range config.updateErrors {

			updateCalls = append(updateCalls, mockClient.EXPECT().
				UpdateOne(
					gomock.AssignableToTypeOf(context.Background()),
					gomock.AssignableToTypeOf(bson.M{}),
					gomock.AssignableToTypeOf(bson.D{})).
				Times(1).
				Return(err))

			// Sorry - if the an error appears then we wont be calling update any more times
			if err != nil {
				break
			}
		}
		gomock.InOrder(updateCalls...)

	} else {
		mockClient.EXPECT().
			StartTransaction(
				gomock.AssignableToTypeOf(context.Background()),
				gomock.AssignableToTypeOf(func() error { return nil })).
			Times(0)

		mockClient.EXPECT().
			UpdateOne(
				gomock.AssignableToTypeOf(context.Background()),
				gomock.AssignableToTypeOf(bson.M{}),
				gomock.AssignableToTypeOf(bson.D{})).
			Times(0)
	}

	return mockClient, ctrl
}

func CreateDummyTransferAccountDate(id int) *model.Account {
	accID := fmt.Sprintf("test%d", id)
	return &model.Account{
		ID:         accID,
		Name:       accID,
		Credential: accID,
		Transactions: []model.Transaction{
			model.Transaction{
				Amount: 1,
				To:     accID,
				From:   "system",
				Date:   time.Now().Format("2006.01.02 15:04:05"),
			},
		},
	}
}

func Test_AccountTransferSucceeds(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       toAccount,
		totalFindCalls: 1,
		unaryTransfer:  false,
		updateErrors:   []error{nil, nil},
	})
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status 200")
	}
}

func Test_AccountTransferFailsOnReceiptientUpdate(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       toAccount,
		totalFindCalls: 1,
		unaryTransfer:  false,
		updateErrors:   []error{nil, errors.New("to account transaction failure")},
	})
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}

func Test_AccountTransferFailsOnCurrentUserUpdate(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       toAccount,
		totalFindCalls: 1,
		unaryTransfer:  false,
		updateErrors:   []error{errors.New("from account transaction failure"), nil},
	})
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}

func Test_AccountTransferFailsOnReceiptientFetch(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      errors.New("fake find error"),
		findData:       toAccount,
		totalFindCalls: 1,
		unaryTransfer:  false,
		updateErrors:   []error{nil, nil},
	})
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected status 404")
	}
}

func Test_AccountTransferFailsOnUnaryTransfer(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       fromAccount,
		totalFindCalls: 1,
		unaryTransfer:  true,
		updateErrors:   []error{nil, nil},
	})
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected status 404")
	}
}

func Test_AccountTransferFailsOnBodyParse(t *testing.T) {
	fromAccount := CreateDummyTransferAccountDate(1)
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       toAccount,
		totalFindCalls: 0,
		unaryTransfer:  false,
		updateErrors:   []error{nil, nil},
	})
	defer ctrl.Finish()

	resp := RunTransferTest(mockClient, fromAccount, []byte{})
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Expected status 400")
	}
}

func Test_AccountTransferFailsAuthAccountNil(t *testing.T) {
	toAccount := CreateDummyTransferAccountDate(2)

	mockClient, ctrl := ConfigureTransferMocks(transferMockConfig{
		t:              t,
		findError:      nil,
		findData:       toAccount,
		totalFindCalls: 0,
		unaryTransfer:  false,
		updateErrors:   []error{nil, nil},
	})
	defer ctrl.Finish()

	resp := RunTransferTest(mockClient, nil, []byte{})
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}
