package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber"
	"github.com/gofiber/logger"

	"dwc.com/lumiere/mongo"
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
	_, err = mongoConfig.Connect("data")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// Set up API
	app := fiber.New()
	api := app.Group("/v1/api", logger.New())
	api.Put("/svcstatus", func(c *fiber.Ctx) {
		c.Status(200).Send("Ok")
	})

	err = app.Listen(config.Port)
	if err != nil {
		log.Printf("Could not start api server on port: %d -> %v", config.Port, err)
	}
}
