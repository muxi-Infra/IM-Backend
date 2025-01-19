package errcode

import (
	"fmt"
)

type Err struct {
	code int
	msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%+v", *e)
}

func NewErr(code int, msg string) *Err {
	return &Err{code: code, msg: msg}
}
func (e *Err) Code() int {
	return e.code
}

func (e *Err) Msg() string {
	return e.msg
}

func (e *Err) WrapError(err error) error {
	return fmt.Errorf("%w:%w", e, err)
}
