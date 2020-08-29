package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber"
	"github.com/gofiber/logger"

	"dwc.com/lumiere/account"
	"dwc.com/lumiere/mongo"
	"dwc.com/lumiere/user"
	"dwc.com/lumiere/utils"
)

type serverConfig struct {
	Port string `env:"PORT,required"`
}

func main() {
	// Parse server configuration
	config := serverConfig{}
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// Parse mongo configuration
	mongoConfig := mongo.MongoConfig{}
	err = env.Parse(&mongoConfig)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// Connect to mongo
	client, collection, err := mongoConfig.Connect("accounts")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	mongoClient := mongo.MongoClient{
		Client: client,
		Conn:   collection,
	}

	// Set up API
	app := fiber.New()
	api := app.Group("/v1/api", logger.New())

	api.Get("/svcstatus", func(c *fiber.Ctx) { c.Status(http.StatusOK).Send("Ok") })

	api.
		Group("/user").
		Post("/register", user.UserRegisterRoute{
			DataAccess: mongoClient,
			Generator:  utils.CodeGenerator{},
		}.Post)

	api.
		Group("/account", user.UserAuthMiddleware{DataAccess: mongoClient}.Auth).
		Get("/balance", account.AccountBalanceRoute{}.GetBalance).
		Get("/transactions", account.AccountBalanceRoute{}.GetTransactions).
		Put("/transfer", account.AccountTransferRoute{DataAccess: mongoClient}.PutTransfer)

	err = app.Listen(config.Port)
	if err != nil {
		log.Printf("Could not start api server on port: %d -> %v", config.Port, err)
	}
}
