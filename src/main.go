package main

import (
	"log"

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
	collection, err := mongoConfig.Connect("accounts")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	client := mongo.MongoClient{
		Conn: collection,
	}

	// Set up API
	app := fiber.New()
	api := app.Group("/v1/api", logger.New())

	api.Get("/svcstatus", func(c *fiber.Ctx) { c.Status(200).Send("Ok") })

	api.Group("/user").
		Post("/register", user.UserRegisterRoute{
			DataAccess: client,
			Generator:  utils.CodeGenerator{},
		}.Post)

	api.Group("/account", user.UserAuthMiddleware{DataAccess: client}.Auth).
		Get("/balance", account.AccountBalanceRoute{DataAccess: client}.Get)

	err = app.Listen(config.Port)
	if err != nil {
		log.Printf("Could not start api server on port: %d -> %v", config.Port, err)
	}
}
