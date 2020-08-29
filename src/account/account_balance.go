package account

import (
	"github.com/gofiber/fiber"

	lMongo "dwc.com/lumiere/mongo"
)

type AccountBalanceRoute struct {
	DataAccess lMongo.IMongoClient
}

func (a AccountBalanceRoute) Get(c *fiber.Ctx) {

	// body := &model.RegisterBody{}
	// if err := c.BodyParser(body); err != nil {
	// 	log.Printf("Could not parse request body: %v", err)
	// 	c.Status(500).Send("Request Invalid")
	// 	return
	// }

	c.Status(200).Send("Authorized")

}
