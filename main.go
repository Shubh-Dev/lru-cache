package main

import (
	"fmt"

	"github.com/Shubh-Dev/lru-cache/cache"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()
	cacheInstance := cache.NewCache(1024)

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Least Recently Used Cache")
	})

	app.Get("/cache/get", func(c *fiber.Ctx) error {
		key := c.Query("key")

		if key == " " {
			return c.Status(fiber.StatusBadRequest).SendString("Key is required")
		}

		value, found := cacheInstance.Get(key)

		if !found {
			return c.Status(fiber.StatusNotFound).SendString("Key not found")
		}
		return c.SendString(fmt.Sprintf("Value for key %s is %v", key, value))
	})

	app.Post("/cache/set", func(c *fiber.Ctx) error {
		return c.SendString("Set cache")
	})

	// Start the server
	app.Listen(":3000")
}
