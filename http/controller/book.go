package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahya077/goChain/service/getBook"
)

func GetBook(ctx *fiber.Ctx) error {
	p := getBook.New().Handle(ctx)

	r := p.GetResponse()

	return ctx.Status(r.Code).JSON(r)
}
