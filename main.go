package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Create a new Fiber instance
    app := fiber.New()

    // Define routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Start the server
    app.Listen(":3000")
}
