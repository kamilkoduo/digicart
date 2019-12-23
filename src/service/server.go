package service

import (
	"github.com/kamilkoduo/digicart/src/api"
)
type CartAPIServer struct{}

func (s CartAPIServer) GetCartType(cartID string) (api.CartType, error) {
	return getCartType(cartID)
}

func (s CartAPIServer) CartExists(cartID string) (bool,error) {
	return cartExists(cartID)
}

func (s CartAPIServer) InitCart(cartID string, cartType api.CartType) (error) {
	return initCart(cartID, cartType)
}

func (s CartAPIServer) MergeCarts(targetCartID, sourceCartID string) error {
	return mergeCarts(targetCartID, sourceCartID)
}
func (s CartAPIServer) GetCart(cartID string) (*api.Cart, error) {
	return getCart(cartID)
}
func (s CartAPIServer) AddCartItem(cartID string, cartItem *api.CartItem) error {
	return addCartItem(cartID, cartItem)
}

func (s CartAPIServer) RemoveCartItem(cartID, cartItemID string) error {
	return removeCartItem(cartID, cartItemID)
}
func (s CartAPIServer) UpdateCartItem(cartID string, cartItem *api.CartItem) error {
	return updateCartItem(cartID, cartItem)
}
