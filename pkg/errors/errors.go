package errors

import "fmt"

type Error interface {
	ErrorCode() ErrorCode
	Error() string
}

type err struct {
	Code ErrorCode `json:"code"`
	Text string    `json:"text"`
}

func New(code ErrorCode, text interface{}, args ...interface{}) Error {
	var errText string
	switch text := text.(type) {
	case error:
		errText = text.Error()
	case string:
		errText = fmt.Sprintf(text, args...)
	}

	return &err{
		Code: code,
		Text: errText,
	}
}

func (e err) ErrorCode() ErrorCode {
	return e.Code
}

func (e err) Error() string {
	return e.Text
}

func Is(err error, target ErrorCode) bool {
	return err.(Error).ErrorCode() == target
}

func Cause(err error) ErrorCode {
	return err.(Error).ErrorCode()
}
