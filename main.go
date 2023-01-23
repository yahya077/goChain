package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yahya077/goChain/http/router"
	"os"
)

func main() {
	app := fiber.New()

	router.Init(app)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
}
