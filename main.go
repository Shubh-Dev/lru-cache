package main

import (
	"fmt"
	"time"

	"github.com/Shubh-Dev/lru-cache/cache"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()
	cacheInstance := cache.NewCache(5)

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		cacheContent := cacheInstance.GetAllCache()

		// Convert the map to a string representation
		var cacheString string
		for key, value := range cacheContent {
			cacheString += fmt.Sprintf("Key: %s, Value: %v\n", key, value)
		}

		if cacheString == "" {
			cacheString = "Cache is empty"
		}

		return c.SendString(cacheString)

	})

	app.Get("/cache/get", func(c *fiber.Ctx) error {
		key := c.Query("key")

		if key == " " {
			return c.Status(fiber.StatusBadRequest).SendString("Key is required")
		}

		value, expirationTime, found := cacheInstance.Get(key)

		if !found {
			return c.Status(fiber.StatusNotFound).SendString("Key not found")
		}
		expirationTimeFormatted := expirationTime.Format("2006-01-02 15:04:05")
		return c.SendString(fmt.Sprintf("Value for key %s is %v, expires at %s", key, value, expirationTimeFormatted))
	})

	app.Post("/cache/set", func(c *fiber.Ctx) error {
		var requestData struct {
			Key        string      `json:"key"`
			Value      interface{} `json:"value"`
			Expiration int         `json:"expiry"`
		}

		if err := c.BodyParser(&requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON format")
		}

		if requestData.Key == " " {
			return c.Status(fiber.StatusBadRequest).SendString("Key is required")
		}

		cacheInstance.Set(requestData.Key, requestData.Value, time.Duration(requestData.Expiration)*time.Second)
		return c.SendString(fmt.Sprintf("Key %s set successfully", requestData.Key))
	})

	// Start the server
	app.Listen(":3000")
}