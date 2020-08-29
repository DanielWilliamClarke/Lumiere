package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
)

type AccountBalanceRoute struct {
	DataAccess lMongo.IMongoClient
}

func (a AccountBalanceRoute) GetBalance(c *fiber.Ctx) {

	account, err := a.getAccount(c)
	if err != nil {
		c.Status(500).Send("Internal server error")
		return
	}

	balance := 0.0
	for _, transaction := range account.Transactions {
		balance = balance + transaction.Amount
	}

	currentTime := time.Now()
	timeFormat := currentTime.Format("2006.01.02 15:04:05")
	balanceReceipt := fmt.Sprintf("%s your current balance at %s is $%.2f", account.Name, timeFormat, balance)

	c.Status(200).Send(balanceReceipt)
}

func (a AccountBalanceRoute) GetTransactions(c *fiber.Ctx) {
	account, err := a.getAccount(c)
	if err != nil {
		c.Status(500).Send("Internal server error")
		return
	}

	c.Status(200).JSON(account.Transactions)
}

func (a AccountBalanceRoute) getAccount(c *fiber.Ctx) (*model.Account, error) {
	authUserID := c.Locals("AuthUser")
	if authUserID == nil {
		reason := "Auth user id is nil"
		log.Printf(reason)
		return nil, errors.New(reason)
	}

	account := &model.Account{}
	err := a.DataAccess.FindOne(context.Background(), bson.M{"_id": authUserID}, account)
	if err != nil {
		log.Printf("User does not exist with given ID: %v", err)
		return nil, err
	}

	return account, nil
}
