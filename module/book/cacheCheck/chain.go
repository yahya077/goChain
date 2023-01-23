package cacheCheck

import "github.com/yahya077/goChain/core/chain"

type Chain struct {
	chain.Chain
}

func (c *Chain) Handle(parameters chain.IParameters) (e error) {
	//you can handle name subject here
	return nil
}

func Handler() chain.IChain {
	return &Chain{chain.NewChain("Elastic Check")}
}
