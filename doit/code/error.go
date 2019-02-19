package code

import (
	"encoding/json"
	"fmt"
)

//根据状态和对应码生成错误信息
func New(status int, code Code) *Error {
	c := new(Error)
	c.Status = status
	c.SetCode(code)
	return c
}

//抽象出来的错误类型
type Error struct {
	Status  int           `json:"-"`
	Code    Code          `json:"code"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func (e Error) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprint(e)
	}
	return string(b)
}

func (e *Error) SetCode(code Code) *Error {
	e.Code = code
	e.Message = parseCodeMessage(code)
	return e
}

func (e *Error) Err(val interface{}, keep ...bool) *Error {
	if len(keep) > 0 && keep[0] == true {
		e.Errors = append(e.Errors, val)
	} else {
		if err, ok := val.(error); ok {
			e.Errors = append(e.Errors, err.Error())
		} else {
			e.Errors = append(e.Errors, val)
		}
	}
	return e
}
