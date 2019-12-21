package carterrors

import "strings"

type ErrType int

const (
	InternalError ErrType = iota
	GeneralError
	CartNotFound
	CartAlreadyExists
	CartItemNotFound
	CartItemAlreadyExists
	InvalidCartID
	InvalidCartItemID
)

var errStr = struct{ m map[ErrType]string }{
	m: map[ErrType]string{
		GeneralError: "General Error : ",
		CartAlreadyExists :"CartAlreadyExists : ",
	},
}

type CartError struct {
	Type ErrType
	msg  string
}

func (ce CartError) Error() string {
	return errStr.m[ce.Type] + ce.msg
}

func (ce CartError) Is(cartError CartError) bool {
	return ce.Type == cartError.Type
}
func (ce CartError) IsType(errType ErrType) bool {
	return ce.Type == errType
}
func New(errType ErrType, msg ...string) error {
	return CartError{
		Type: errType,
		msg:  strings.Join(msg,","),
	}
}