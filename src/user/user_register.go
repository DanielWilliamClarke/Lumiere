package user

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber"

	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
	"dwc.com/lumiere/utils"
)

type UserRegisterRoute struct {
	DataAccess lMongo.IMongoClient
	Generator  utils.ICodeGenerator
}

type registerOutput struct {
	UserID     string `json:"userID"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

type RegisterBody struct {
	Username string  `json:"username",form:"username"`
	Cash     float64 `json:"cash",form:"cash"`
}

func (a UserRegisterRoute) Post(c *fiber.Ctx) {

	body := &RegisterBody{}
	if err := c.BodyParser(body); err != nil {
		log.Printf("Could not parse request body: %v", err)
		c.Status(500).Send("Request Invalid")
		return
	}

	account := &model.Account{}
	err := a.DataAccess.FindOne(context.Background(), bson.M{"name": body.Username}, account)
	if err == nil {
		log.Printf("User already exists: %v", err)
		c.Status(500).Send("User already exists with that user name")
		return
	}

	auth, err := a.Generator.Generate(16)
	if err != nil {
		log.Printf("Could not generate user code: %v", err)
		c.Status(500).Send("Could not register user")
		return
	}

	userID, err := a.Generator.Generate(5)
	if err != nil {
		log.Printf("Could not generate user ID: %v", err)
		c.Status(500).Send("Could not register user")
		return
	}

	data, err := bson.Marshal(&model.Account{
		ID:         userID,
		Name:       body.Username,
		Credential: auth,
		Transactions: []model.Transaction{
			model.Transaction{
				Amount: body.Cash,
				To:     userID,
				From:   "system",
				Date:   time.Now().Format("2006.01.02 15:04:05"),
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
		UserID:     userID,
		Username:   body.Username,
		Credential: auth,
	})
}
