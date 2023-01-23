package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahya077/goChain/http/controller"
)

func Book(app *fiber.App) {
	app.Get("/", controller.GetBook)
}
