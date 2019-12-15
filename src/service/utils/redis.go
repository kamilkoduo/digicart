package utils

import (
	"github.com/kamilkoduo/digicart/src/api"
)

func CartKey(cart *api.Cart) string {
	return "cart:" + cart.CartID
}
func CartInfoMap(cart *api.Cart) (map[string]interface{}) {
	return map[string]interface{}{
		"user_id":  cart.UserID,
		"guest_id": cart.GuestID,
	}
}
func CartItemsKey(cart *api.Cart) string {
	return CartKey(cart) + ":items"
}
