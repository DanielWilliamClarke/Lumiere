package user

import (
	"context"
	"log"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
)

type UserAuthMiddleware struct {
	DataAccess lMongo.IMongoClient
}

type UserCodeBody struct {
	UserCode string `json:"userCode",form:"userCode"`
}

func (u UserAuthMiddleware) Auth(c *fiber.Ctx) {

	body := &UserCodeBody{}
	if err := c.BodyParser(body); err != nil {
		log.Printf("Could not parse request user code: %v", err)
		c.Status(403).Send("User not Authorized")
		return
	}

	account := &model.Account{}
	err := u.DataAccess.FindOne(context.Background(), bson.M{"credential": body.UserCode}, account)
	if err != nil {
		log.Printf("User does not exist with given usercode: %v", err)
		c.Status(403).Send("User not Authorized")
		return
	}

	c.Locals("id", account.M_ID)
	c.Next()
}
