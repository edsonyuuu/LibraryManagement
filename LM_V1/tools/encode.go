package tools

import "time"

const (
	NotLogin     = 10086
	UserError    = 10010
	NotRegister  = 10011
	OK           = 2001
	DoErr        = 2002
	NotFound     = 2003
	Failed       = 2004
	ContinueTime = 30 * 24 * time.Hour
)

type HttpCode struct {
	Code    int         `json:"code"`
	Message string      `json:"name"`
	Data    interface{} `json:"data"`
}
