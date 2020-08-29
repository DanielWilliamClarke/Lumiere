package account

import (
	"context"
	"log"

	"github.com/gofiber/fiber"

	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
	"dwc.com/lumiere/utils"
)

type AccountRegisterRoute struct {
	DataAccess lMongo.IMongoClient
}

type registerOutput struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	UserCode string `json:"usercode"`
}

func (a AccountRegisterRoute) Post(c *fiber.Ctx) {

	body := &model.RegisterBody{}
	if err := c.BodyParser(body); err != nil {
		log.Printf("Could not parse request body: %v", err)
		c.Status(500).Send("Request Invalid")
		return
	}

	account := &model.Account{}
	err := a.DataAccess.FindOne(context.Background(), bson.D{{"name", body.Username}}, account)
	if err == nil {
		log.Printf("User already exists: %v", err)
		c.Status(500).Send("User already exists with that user name")
		return
	}

	userCode, err := utils.GenerateUID(10)
	if err != nil {
		log.Printf("Could not generate user code: %v", err)
		c.Status(500).Send("Could not register user")
		return
	}

	userID, err := utils.GenerateUID(5)
	if err != nil {
		log.Printf("Could not generate user ID: %v", err)
		c.Status(500).Send("Could not register user")
		return
	}

	data, err := bson.Marshal(&model.Account{
		ID:       userID,
		Name:     body.Username,
		UserCode: userCode,
		Transactions: []model.Transaction{
			model.Transaction{
				Amount: body.Cash,
				To:     userID,
				From:   "system",
			},
		},
	})
	if err != nil {
		log.Printf("User malformed: %v", err)
		c.Status(500).Send("Could not generate user")
		return
	}

	_, err = a.DataAccess.InsertOne(context.Background(), data)
	if err != nil {
		log.Printf("Could not insert user: %v", err)
		c.Status(500).Send("Could not insert user")
		return
	}

	c.Status(200).JSON(registerOutput{
		UserID:   userID,
		Username: body.Username,
		UserCode: userCode,
	})
}
