package test

import (
	"bytes"
	"context"
	"encoding/json"
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

func ConfigureTransferMocks(t *testing.T, findData interface{}, findError error, totalCalls int) (*mock_mongo.MockIMongoClient, *gomock.Controller) {

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

	mockClient.EXPECT().
		UpdateOne(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(bson.M{}),
			gomock.AssignableToTypeOf(bson.D{})).
		Times(2).
		Return(nil)

	mockClient.EXPECT().
		StartTransaction(
			gomock.AssignableToTypeOf(context.Background()),
			gomock.AssignableToTypeOf(func() error { return nil })).
		Times(1).
		DoAndReturn(func(ctx context.Context, callback func() error) error {
			return callback()
		})

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

	mockClient, ctrl := ConfigureTransferMocks(t, toAccount, nil, 1)
	defer ctrl.Finish()

	body, err := json.Marshal(user.RegisterBody{
		Username: "test",
		Cash:     100,
	})
	if err != nil {
		t.Errorf("Could not marshal body: %v", err)
	}

	resp := RunTransferTest(mockClient, fromAccount, body)
	if resp.StatusCode != 200 {
		t.Error("Expected status 200")
	}
}
