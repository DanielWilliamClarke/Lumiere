package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"

	"dwc.com/lumiere/utils"
)

func RunRequestDurationTest() *http.Response {

	// Test fiber routing logic with
	// https://docs.gofiber.io/app#test

	// Create request
	req := httptest.NewRequest("GET", "http://localhost:5000", nil)
	req.Header.Set("Content-Type", "application/json")

	// Setup Test router
	app := fiber.New()
	app.Use(utils.RequestDurationMonitor()).
		Get("/", func(c *fiber.Ctx) {
			c.Status(http.StatusOK).Send("Authorized")
		})

	// Run Test
	resp, _ := app.Test(req)

	return resp
}

func Test_RequestDurationMiddlewareSucceeds(t *testing.T) {
	resp := RunRequestDurationTest()
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status 200")
	}
}
