package router

import "github.com/gofiber/fiber/v2"

func Init(app *fiber.App) {
	Book(app)
}
