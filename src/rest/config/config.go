package config

import (
	"fmt"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service"
	service_db "github.com/kamilkoduo/digicart/src/service/db/service"
	"os"
)

// APIServer ...
var APIServer api.CartAPI = func() api.CartAPI {
	s := &service.CartAPIServer{}
	dbs := &service_db.CartDBRedisServer{}
	fmt.Printf("DDD %v %T", s.GetUnexportedCartDBAPIServer(), s.GetUnexportedCartDBAPIServer())
	s.RegisterCartDBAPIServer(dbs)
	fmt.Printf("DDD %v %T", s.GetUnexportedCartDBAPIServer(), s.GetUnexportedCartDBAPIServer())
	return s
}()

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
