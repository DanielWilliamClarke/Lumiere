package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/gofiber/fiber"
	"github.com/ulule/deepcopier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"dwc.com/lumiere/account"
	"dwc.com/lumiere/account/model"
	"dwc.com/lumiere/mocks"
	mock_mongo "dwc.com/lumiere/mocks"
	lMongo "dwc.com/lumiere/mongo"
)

func RunBalanceTest(mockClient lMongo.IMongoClient, endpoint string, authUser interface{}) *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:5000/%s", endpoint), nil)

	// Setup Test router
	app := fiber.New()
	app.
		Use(func(c *fiber.Ctx) {
			c.Locals("AuthUser", authUser)
			c.Next()
		}).
		Get("/b", account.AccountBalanceRoute{DataAccess: mockClient}.GetBalance).
		Get("/t", account.AccountBalanceRoute{DataAccess: mockClient}.GetTransactions)

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

func ConfigureBalanceMocks(t *testing.T, findData *model.Account, findError error, totalCalls int) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockIMongoClient(ctrl)

	mockClient.EXPECT().
		FindOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(&model.Account{})).
		Times(totalCalls).
		DoAndReturn(func(ctx context.Context, filter interface{}, data interface{}) error {
			if findData != nil {
				deepcopier.Copy(findData).To(data)
			}
			return findError
		})

	return mockClient, ctrl
}

func CreateDummyAccountDate() *model.Account {
	return &model.Account{
		ID:         "test",
		Name:       "test",
		Credential: "tst",
		Transactions: []model.Transaction{
			model.Transaction{
				Amount: 1,
				To:     "test",
				From:   "system",
				Date:   time.Now().Format("2006.01.02 15:04:05"),
			},
			model.Transaction{
				Amount: 1,
				To:     "test",
				From:   "system",
				Date:   time.Now().Format("2006.01.02 15:04:05"),
			},
		},
	}
}

func Test_AccountBalanceCanBeReturned(t *testing.T) {

	account := CreateDummyAccountDate()

	mockClient, ctrl := ConfigureBalanceMocks(t, account, nil, 1)
	defer ctrl.Finish()

	resp := RunBalanceTest(mockClient, "b", &primitive.ObjectID{})
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}

func Test_AccountBalanceFindFailure(t *testing.T) {

	account := CreateDummyAccountDate()

	mockClient, ctrl := ConfigureBalanceMocks(t, account, errors.New("fake find failure"), 1)
	defer ctrl.Finish()

	resp := RunBalanceTest(mockClient, "b", &primitive.ObjectID{})
	if resp.StatusCode != 500 {
		t.Error("Expected status 500")
	}
}

func Test_AccountAuthUserFailure(t *testing.T) {

	account := CreateDummyAccountDate()

	mockClient, ctrl := ConfigureBalanceMocks(t, account, nil, 0)
	defer ctrl.Finish()

	resp := RunBalanceTest(mockClient, "b", nil)
	if resp.StatusCode != 500 {
		t.Error("Expected status 500")
	}
}

func Test_AccountTransactionsCanBeReturned(t *testing.T) {

	account := CreateDummyAccountDate()

	mockClient, ctrl := ConfigureBalanceMocks(t, account, nil, 1)
	defer ctrl.Finish()

	resp := RunBalanceTest(mockClient, "t", &primitive.ObjectID{})
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}

func Test_AccountTransactionsFindFailure(t *testing.T) {

	account := CreateDummyAccountDate()

	mockClient, ctrl := ConfigureBalanceMocks(t, account, errors.New("fake find failure"), 1)
	defer ctrl.Finish()

	resp := RunBalanceTest(mockClient, "t", &primitive.ObjectID{})
	if resp.StatusCode != 500 {
		t.Error("Expected status 500")
	}
}
