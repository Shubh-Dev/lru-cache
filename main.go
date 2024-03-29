package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Shubh-Dev/lru-cache/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	cacheInstance := cache.NewCache(5)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		cacheContent := cacheInstance.GetAllCache()
		return c.JSON(cacheContent)

	})

	app.Get("/cache/get", func(c *fiber.Ctx) error {
		key := c.Query("key")

		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Key is required",
			})
		}

		value, expirationTime, found := cacheInstance.Get(key)

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Key not found",
			})
		}

		expirationTimeFormatted := expirationTime.Format("2006-01-02 15:04:05")
		response := fiber.Map{
			"key":    key,
			"value":  value,
			"expiry": expirationTimeFormatted,
		}

		return c.JSON(response)
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := app.Listen(":" + port); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

}
