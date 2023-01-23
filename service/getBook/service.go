package getBook

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahya077/goChain/core/chain"
	"github.com/yahya077/goChain/module/book/chainParameters"
	"github.com/yahya077/goChain/module/book/manager"
)

type Service struct {
}

func (s Service) Handle(ctx *fiber.Ctx) chain.IParameters {
	p := chainParameters.NewParameters(ctx)

	// you can add parameters here (eq: request payload)
	// ...

	return chain.Handle(manager.NewGetBookManager(p))
}

func New() Service {
	return Service{}
}
