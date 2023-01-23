package manager

import (
	"github.com/yahya077/goChain/core/chain"
	"github.com/yahya077/goChain/module/book/cacheCheck"
	"github.com/yahya077/goChain/module/book/elasticCheck"
	"github.com/yahya077/goChain/module/book/response"
	"github.com/yahya077/goChain/module/book/visibilityCheck"
)

type GetBook struct {
	parameters chain.IParameters
}

func (m GetBook) BuildChains() chain.IChain {
	visibilityCheckHandler := visibilityCheck.Handler()
	elasticCheckHandler := elasticCheck.Handler()
	cacheCheckHandler := cacheCheck.Handler()
	responseHandler := response.Handler()

	visibilityCheckHandler.SetSuccessHandler(cacheCheckHandler)
	cacheCheckHandler.SetSuccessHandler(responseHandler)
	cacheCheckHandler.SetSkipHandler(elasticCheckHandler)
	elasticCheckHandler.SetSuccessHandler(responseHandler)

	return visibilityCheckHandler
}

func (m GetBook) GetParameters() chain.IParameters {
	return m.parameters
}

func NewGetBookManager(parameters chain.IParameters) chain.IChainManager {
	return GetBook{parameters}
}
