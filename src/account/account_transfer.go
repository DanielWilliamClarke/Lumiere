package account

import (
	"context"
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

	authedAccount := c.Locals("AuthedAccount").(*model.Account)
	if authedAccount == nil {
		log.Printf("Auth user id is nil")
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	body := &TransferBody{}
	if err := c.BodyParser(body); err != nil {
		log.Printf("Could not parse request body: %v", err)
		c.Status(http.StatusBadRequest).Send("Request Invalid")
		return
	}

	ctx := context.Background()

	toAccount := &model.Account{}
	err := a.DataAccess.FindOne(ctx, bson.M{"name": body.To}, toAccount)
	if err != nil {
		log.Printf("Transfer receiptient does not exist: %v", err)
		c.Status(http.StatusNotFound).Send("Internal server error")
		return
	}

	if toAccount.ID == authedAccount.ID {
		log.Printf("Transfer receiptient same as sender: %v", err)
		c.Status(http.StatusBadRequest).Send("Cannot transfer to self")
		return
	}

	authedAccount.Transactions = append(authedAccount.Transactions, model.Transaction{
		Amount: -body.Amount,
		From:   authedAccount.ID,
		To:     toAccount.ID,
		Date:   time.Now().Format("2006.01.02 15:04:05"),
	})
	toAccount.Transactions = append(toAccount.Transactions, model.Transaction{
		Amount: body.Amount,
		From:   authedAccount.ID,
		To:     toAccount.ID,
		Date:   time.Now().Format("2006.01.02 15:04:05"),
	})

	// Perform update in transaction so we can rollback all changes if any error occurs here
	err = a.DataAccess.StartTransaction(ctx, func() error {
		err = a.DataAccess.UpdateOne(ctx, bson.M{"_id": authedAccount.M_ID}, bson.D{{"$set", bson.M{"transactions": authedAccount.Transactions}}})
		if err != nil {
			log.Printf("Unable to update [from] account transactions: %v", err)
			return err
		}
		err = a.DataAccess.UpdateOne(ctx, bson.M{"_id": toAccount.M_ID}, bson.D{{"$set", bson.M{"transactions": toAccount.Transactions}}})
		if err != nil {
			log.Printf("Unable to update [to] account transactions: %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Account transfer failed: %v", err)
		c.Status(http.StatusInternalServerError).Send("Internal server error")
		return
	}

	timeFormat := time.Now().Format("2006.01.02 15:04:05")
	transferReceipt := fmt.Sprintf("%s your transfer of $%.2f at %s to %s is complete", authedAccount.Name, body.Amount, timeFormat, toAccount.Name)
	c.Status(http.StatusOK).Send(transferReceipt)
}
