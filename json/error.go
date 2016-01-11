package json

import (
	"errors"
)

var ErrNullResult = errors.New("result is null")

type Error struct {
	Code    int         `json:"Code"`
	Message string      `json:"mesage"`
	Data    interface{} `json:"data"`
}

func (e *Error) Error() string {
	return e.Message
}
