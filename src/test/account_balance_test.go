package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber"

	"dwc.com/lumiere/account"
	"dwc.com/lumiere/account/model"
)

func RunBalanceTest(endpoint string, authedAccount interface{}) *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:5000/%s", endpoint), nil)

	// Setup Test router
	app := fiber.New()
	app.
		Use(func(c *fiber.Ctx) {
			c.Locals("AuthedAccount", authedAccount)
			c.Next()
		}).
		Get("/b", account.AccountBalanceRoute{}.GetBalance).
		Get("/t", account.AccountBalanceRoute{}.GetTransactions)

	// Run Test
	resp, _ := app.Test(req)

	return resp
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
	resp := RunBalanceTest("b", account)
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status 200")
	}
}

func Test_AccountBalanceAuthUserFailure(t *testing.T) {
	resp := RunBalanceTest("b", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}

func Test_AccountTransactionsCanBeReturned(t *testing.T) {
	account := CreateDummyAccountDate()
	resp := RunBalanceTest("t", account)
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status 200")
	}
}

func Test_AccountTransactionsAuthUserFailure(t *testing.T) {
	resp := RunBalanceTest("t", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500")
	}
}
