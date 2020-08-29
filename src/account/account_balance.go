package account

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber"

	"dwc.com/lumiere/account/model"
)

type AccountBalanceRoute struct{}

func (a AccountBalanceRoute) GetBalance(c *fiber.Ctx) {

	account, err := a.getAccount(c)
	if err != nil {
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	balance := 0.0
	for _, transaction := range account.Transactions {
		balance = balance + transaction.Amount
	}

	timeFormat := time.Now().Format("2006.01.02 15:04:05")
	balanceReceipt := fmt.Sprintf("%s your current balance at %s is $%.2f", account.Name, timeFormat, balance)

	c.Status(http.StatusOK).Send(balanceReceipt)
}

func (a AccountBalanceRoute) GetTransactions(c *fiber.Ctx) {
	account, err := a.getAccount(c)
	if err != nil {
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	c.Status(http.StatusOK).JSON(account.Transactions)
}

func (a AccountBalanceRoute) getAccount(c *fiber.Ctx) (*model.Account, error) {
	authedAccount := c.Locals("AuthedAccount")
	if authedAccount == nil {
		reason := "Authed account is nil"
		log.Printf(reason)
		return nil, errors.New(reason)
	}

	return authedAccount.(*model.Account), nil
}
