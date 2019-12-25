package carterrors

import (
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

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
	Unauthenticated
	AccessDenied
)

var title = struct{ m map[ErrType]string }{
	m: map[ErrType]string{
		InternalError:         "InternalError",
		GeneralError:          "GeneralError",
		CartNotFound:          "CartNotFound",
		CartAlreadyExists:     "CartAlreadyExists",
		CartItemNotFound:      "CartItemNotFound",
		CartItemAlreadyExists: "CartItemAlreadyExists",
		InvalidCartID:         "InvalidCartID",
		InvalidCartItemID:     "InvalidCartItemID",
		Unauthenticated:       "Unauthenticated",
		AccessDenied:          "AccessDenied",
	},
}

func (ce CartError) Title() string {
	return title.m[ce.Type()]
}
func (ce CartError) AppCode() int {
	return int(ce.Type())
}
func (ce CartError) AppCodeText() string {
	return strconv.Itoa(ce.AppCode())
}

var statusCode = struct{ m map[ErrType]int }{
	m: map[ErrType]int{
		InternalError:         http.StatusInternalServerError,
		GeneralError:          http.StatusInternalServerError,
		CartNotFound:          http.StatusNotFound,
		CartAlreadyExists:     http.StatusConflict,
		CartItemNotFound:      http.StatusNotFound,
		CartItemAlreadyExists: http.StatusConflict,
		InvalidCartID:         http.StatusBadRequest,
		InvalidCartItemID:     http.StatusBadRequest,
		Unauthenticated:       http.StatusUnauthorized,
		AccessDenied:          http.StatusForbidden,
	},
}

func (ce CartError) StatusCode() int {
	return statusCode.m[ce.Type()]
}
func (ce CartError) StatusCodeText() string {
	return http.StatusText(ce.StatusCode())
}

type CartError struct {
	id     string
	ceType ErrType
	msg    string
}

func (ce CartError) ID() string {
	return ce.id
}
func (ce CartError) Type() ErrType {
	return ce.ceType
}
func (ce CartError) Message() string {
	return ce.msg
}

func (ce CartError) Error() string {
	return strings.Join([]string{ce.Title(), ce.Message()}, ":")
}

func (ce CartError) Is(cartError CartError) bool {
	return ce.IsType(cartError.Type())
}

func (ce CartError) IsType(errType ErrType) bool {
	return ce.Type() == errType
}

func New(errType ErrType, msg ...string) CartError {
	return CartError{
		id:     uuid.New().String(),
		ceType: errType,
		msg:    strings.Join(msg, ","),
	}
}

func (ce CartError) JSONObject() *jsonapi.ErrorObject {
	eo := &jsonapi.ErrorObject{
		ID:     ce.ID(),
		Title:  ce.Title(),
		Detail: ce.Message(),
		Status: ce.StatusCodeText(),
		Code:   ce.AppCodeText(),
	}
	return eo
}
