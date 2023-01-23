package chain

import (
	"fmt"
	"strings"
	"time"
)

type IChain interface {
	IsProcessable(parameters IParameters) bool
	GetName() string
	SetSuccessHandler(chain IChain)
	getSuccessHandler() IChain
	SetErrorHandler(chain IChain)
	getErrorHandler() IChain
	SetSkipHandler(chain IChain)
	getSkipHandler() IChain
	Handle(parameters IParameters) error
	init(handler IChain, parameters IParameters) IParameters
	ThrowError(data any, message string, code int, parameters IParameters) (err error)
	SetResponse(data any, message string, code int, parameters IParameters)
}

// Chain is the base chain struct
type Chain struct {
	Name           string
	SuccessHandler IChain
	ErrorHandler   IChain
	SkipHandler    IChain
	Logger         IChainLogger
}

func NewChain(name string) Chain {
	return Chain{
		Name:   name,
		Logger: NewChainLogger(NewLogger()),
	}
}

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

func (c *Chain) IsProcessable(parameters IParameters) bool {
	return true
}

// Handle you should override this function
func (c *Chain) Handle(parameters IParameters) error {
	fmt.Println("On Handle... You should override")
	return nil
}

func (c *Chain) init(handler IChain, parameters IParameters) IParameters {
	logStartAt := time.Now()
	c.Logger.LogStart(c.Name)
	if handler.IsProcessable(parameters) {
		if e := handler.Handle(parameters); e != nil {
			c.Logger.LogFinishWithError(c.Name, logStartAt, e)

			if next := c.getErrorHandler(); next != nil {
				return next.init(next, parameters)
			}

			return parameters
		}
		c.Logger.LogFinish(c.Name, logStartAt)
	} else {
		c.Logger.LogSkipped(c.Name, logStartAt)
		if next := c.getSkipHandler(); next != nil {
			return next.init(next, parameters)
		}
	}

	if parameters.Completed() {
		return parameters
	}

	if next := c.getSuccessHandler(); next != nil {
		return next.init(next, parameters)
	}

	return parameters
}

func (c *Chain) GetName() string {
	return c.Name
}

func (c *Chain) getSuccessHandler() IChain {
	return c.SuccessHandler
}

func (c *Chain) SetSuccessHandler(chain IChain) {
	c.SuccessHandler = chain
}

func (c *Chain) getErrorHandler() IChain {
	return c.ErrorHandler
}

func (c *Chain) SetErrorHandler(chain IChain) {
	c.ErrorHandler = chain
}

func (c *Chain) getSkipHandler() IChain {
	return c.SkipHandler
}

func (c *Chain) SetSkipHandler(chain IChain) {
	c.SkipHandler = chain
}

func (c *Chain) ThrowError(data any, message string, code int, parameters IParameters) (err error) {
	c.SetResponse(data, message, code, parameters)
	return
}

func (c *Chain) SetResponse(data any, message string, code int, parameters IParameters) {
	r := Response{
		Code:    code,
		Data:    data,
		Message: message,
	}

	parameters.SetResponse(r)
}

/***** Parameters *****/

type Parameters struct {
	IsCompleted bool
	Response    Response
	Extra       map[string]interface{}
}

func (p *Parameters) GetExtra(key string) any {
	return p.Extra[key]
}

func (p *Parameters) AddExtra(key string, v any) {
	p.Extra[key] = v
}

func (p *Parameters) Completed() bool {
	return p.IsCompleted
}

func (p *Parameters) SetResponse(r Response) {
	p.IsCompleted = true
	p.Response = r
}

func (p *Parameters) GetResponse() Response {
	return p.Response
}

type IParameters interface {
	GetExtra(key string) any
	AddExtra(key string, v any)
	SetResponse(r Response)
	GetResponse() Response
	Completed() bool
}

func NewParameters() Parameters {
	return Parameters{Extra: map[string]interface{}{}}
}

/***** Logger ******/

type LogType string

const (
	Success LogType = "SUCCESS"
	Info    LogType = "INFO"
	Error   LogType = "ERROR"
	Warning LogType = "WARNING"
)

type ILogger interface {
	Log(logType LogType, message string, v ...interface{})
}

type Logger struct {
}

func (l Logger) Log(logType LogType, message string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("%s: %s", logType, message))
}

func NewLogger() Logger {
	return Logger{}
}

type IChainLogger interface {
	LogStart(on string)
	LogFinish(on string, startAt time.Time)
	LogSkipped(on string, startAt time.Time)
	LogFinishWithError(on string, startAt time.Time, err error)
}

type ChainLogger struct {
	Logger ILogger
}

func NewChainLogger(logger ILogger) IChainLogger {
	return ChainLogger{logger}
}

func (l ChainLogger) LogStart(on string) {
	l.Logger.Log(Info, fmt.Sprintf(": (%s) Start", on))
}

func (l ChainLogger) LogFinish(on string, startAt time.Time) {
	l.Logger.Log(Info, fmt.Sprintf(": (%s) End, took: %s", on, normalizeDuration(time.Since(startAt))))
}

func (l ChainLogger) LogSkipped(on string, startAt time.Time) {
	l.Logger.Log(Info, fmt.Sprintf(": (%s) End (Skipped), took: %s", on, normalizeDuration(time.Since(startAt))))
}

func (l ChainLogger) LogFinishWithError(on string, startAt time.Time, err error) {
	l.Logger.Log(Error, fmt.Sprintf(": (%s) End, took: %s.\nWith Error; %s", on, normalizeDuration(time.Since(startAt)), err), err)
}

type IChainManager interface {
	BuildChains() IChain
	GetParameters() IParameters
}

/***** Manager ******/

func Handle(chainManager IChainManager) IParameters {
	chains := chainManager.BuildChains()

	return chains.init(chains, chainManager.GetParameters())
}

/***** Helper func ******/

func normalizeDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
