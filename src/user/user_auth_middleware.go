package user

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"dwc.com/lumiere/account/model"
	lMongo "dwc.com/lumiere/mongo"
)

type UserAuthMiddleware struct {
	DataAccess lMongo.IMongoClient
}

func (u UserAuthMiddleware) Auth(c *fiber.Ctx) {

	auth := c.Get("Authorization")
	if len(auth) == 0 {
		log.Printf("User authorization not present in header")
		c.Status(http.StatusUnauthorized).Send("User not Authorized")
		return
	}

	account := &model.Account{}
	err := u.DataAccess.FindOne(context.Background(), bson.M{"credential": auth}, account)
	if err != nil {
		log.Printf("User does not exist with given usercode: %v", err)
		c.Status(http.StatusUnauthorized).Send("User not Authorized")
		return
	}

	c.Locals("AuthedAccount", account)
	c.Next()
}
