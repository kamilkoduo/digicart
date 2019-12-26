package config

import (
	"github.com/kamilkoduo/digicart/src/service"
	"os"
)

// APIServer ...
var APIServer = service.CartAPIServer{}


/* app consts*/
const defaultAppAddress = "0.0.0.0:8080"

// AppAddress ...
var AppAddress = func() string {
	val, found := os.LookupEnv("APP_ADDRESS")
	if !found {
		val = defaultAppAddress
	}
	return val
}()

// headers
const (
	HeaderUserID      = "x-user-id"
	HeaderGuestID     = "x-guest-id"
	HeaderContentType = "content-type"
)

// paths
const (
	PathAPIv1My      = "/api/v1/cart-api/my"
	PathPostfixItems = "/items"
)

// keys
const (
	KeyCartID      = "cartID"
	KeyCartType    = "cartType"
	KeyMergeCartID = "mergeCartID"
	KeyID          = "id"
	KeyCartItem    = "cartItem"
	KeyCartItemID  = "cartItemID"
)
