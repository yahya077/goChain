package chainParameters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahya077/goChain/core/chain"
)

const MyExtra = "my_example_extra"

type Parameters struct {
	chain.Parameters
}

func NewParameters(ctx *fiber.Ctx) *Parameters {
	p := &Parameters{chain.NewParameters()}

	// add some extra to parameters if you needed
	p.AddExtra(MyExtra, "this is an example")

	return p
}

type ParametersHelper struct {
	// you can create your base struct and init here as abstract
}

func (p ParametersHelper) GetExample(parameters chain.IParameters) string {
	return parameters.GetExtra(MyExtra).(string)
}

func Helper() ParametersHelper {
	return ParametersHelper{}
}
