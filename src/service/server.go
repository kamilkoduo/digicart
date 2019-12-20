package service

import (
	"github.com/kamilkoduo/digicart/src/api"
)
type CartApiServer struct{}

func (s CartApiServer) GetCartType(cartID string) (api.CartType, error) {
	return getCartType(cartID)
}

func (s CartApiServer) CartExists(cartID string) (bool,error) {
	return cartExists(cartID)
}

func (s CartApiServer) InitCart(cartID string, cartType api.CartType) (error) {
	return initCart(cartID, cartType)
}

func (s CartApiServer) MergeCarts(targetCartID string, sourceCartID string) error {
	return mergeCarts(targetCartID, sourceCartID)
}
func (s CartApiServer) GetCart(cartID string) (*api.Cart, error) {
	return getCart(cartID)
}
func (s CartApiServer) AddCartItem(cartID string, cartItem *api.CartItem) error {
	return addCartItem(cartID, cartItem)
}

func (s CartApiServer) RemoveCartItem(cartID string, cartItemID string) error {
	return removeCartItem(cartID, cartItemID)
}
func (s CartApiServer) UpdateCartItem(cartID string, cartItem *api.CartItem) error {
	return updateCartItem(cartID, cartItem)
}
