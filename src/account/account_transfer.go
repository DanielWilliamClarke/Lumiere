package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
)

type AccountTransferRoute struct {
	DataAccess lMongo.IMongoClient
}

type TransferBody struct {
	To     string  `json:"to",form:"to"`
	Amount float64 `json:"amount",form:"amount"`
}

func (a AccountTransferRoute) PutTransfer(c *fiber.Ctx) {

	ctx := context.Background()

	// Retrieve authed account from upstream
	authedAccount, err := a.getAccount(c)
	if err != nil {
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	// Parse request body
	body := &TransferBody{}
	if err := c.BodyParser(body); err != nil {
		log.Printf("Could not parse request body: %v", err)
		c.Status(http.StatusBadRequest).Send("Request Invalid")
		return
	}

	// Retrieve receiptent from storage
	receiptient, err := a.getReceiptient(ctx, body.To, authedAccount)
	if err != nil {
		c.Status(http.StatusNotFound).Send(err.Error())
		return
	}

	// Perform update in transaction so we can rollback all changes if any error occurs here
	err = a.safeTransfer(ctx, authedAccount, receiptient, body.Amount)
	if err != nil {
		log.Printf("Account transfer failed: %v", err)
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	// Render success output
	timeFormat := time.Now().Format("2006.01.02 15:04:05")
	transferReceipt := fmt.Sprintf("%s your transfer of $%.2f at %s to %s is complete", authedAccount.Name, body.Amount, timeFormat, receiptient.Name)
	c.Status(http.StatusOK).Send(transferReceipt)
}

func (a AccountTransferRoute) getAccount(c *fiber.Ctx) (*model.Account, error) {
	authedAccount := c.Locals("AuthedAccount")
	if authedAccount == nil {
		reason := "Authed account is nil"
		log.Printf(reason)
		return nil, errors.New(reason)
	}

	return authedAccount.(*model.Account), nil
}

func (a AccountTransferRoute) getReceiptient(ctx context.Context, receiptientName string, authedAccount *model.Account) (*model.Account, error) {
	// Retrive to receiptient account
	receiptient := &model.Account{}
	err := a.DataAccess.FindOne(ctx, bson.M{"name": receiptientName}, receiptient)
	if err != nil {
		log.Printf("Transfer receiptient does not exist: %v", err)
		return nil, errors.New("Transfer receiptient does not exist")
	}

	// Guard agaisnt unary transfer
	if receiptient.ID == authedAccount.ID {
		log.Printf("Transfer receiptient same as sender: %v", err)
		return nil, errors.New("Cannot transfer to self")
	}

	return receiptient, nil
}

func (a AccountTransferRoute) safeTransfer(ctx context.Context, authedAccount *model.Account, receiptient *model.Account, amount float64) error {
	return a.DataAccess.StartTransaction(ctx, func() error {

		err := a.DataAccess.UpdateOne(ctx,
			bson.M{"_id": authedAccount.M_ID},
			a.createTransactionUpdate(authedAccount, receiptient, -amount))
		if err != nil {
			log.Printf("Unable to update [from] account transactions: %v", err)
			return err
		}

		err = a.DataAccess.UpdateOne(ctx,
			bson.M{"_id": receiptient.M_ID},
			a.createTransactionUpdate(authedAccount, receiptient, amount))
		if err != nil {
			log.Printf("Unable to update [to] account transactions: %v", err)
			return err
		}

		return nil
	})
}

func (a AccountTransferRoute) createTransactionUpdate(authedAccount *model.Account, receiptient *model.Account, amount float64) bson.D {
	return bson.D{
		{"$addToSet",
			bson.M{"transactions": model.Transaction{
				Amount: amount,
				From:   authedAccount.ID,
				To:     receiptient.ID,
				Date:   time.Now().Format("2006.01.02 15:04:05"),
			}},
		}}
}
