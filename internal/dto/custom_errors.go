package dto

import (
	"errors"
	"fmt"
)

var (
	NotFoundError = errors.New("not found")
)

// ErrorResponse a wrapper for error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}

func (e ErrorResponse) String() string {
	return fmt.Sprintf("Error response status:%d, Error:%v, Message:%s", e.Status, e.Error, e.Message)
}
