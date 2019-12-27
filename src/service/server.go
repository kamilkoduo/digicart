package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	db_api "github.com/kamilkoduo/digicart/src/service/db/api"
	"log"
)

// CartAPIServer ...
type CartAPIServer struct {
	cartDBAPI db_api.CartDBAPI
}

func (s *CartAPIServer) RegisterCartDBAPIServer(dbs db_api.CartDBAPI) {
	log.Printf("Registration of cart DB API Server of type: %T\n", dbs)
	s.cartDBAPI = dbs
}

func (s *CartAPIServer) GetCartType(cartID string) (api.CartType, error) {
	return s.getCartType(cartID)
}

func (s *CartAPIServer) CartExists(cartID string) (bool, error) {
	return s.cartExists(cartID)
}

func (s *CartAPIServer) InitCart(cartID string, cartType api.CartType) (error) {
	return s.initCart(cartID, cartType)
}

func (s *CartAPIServer) MergeCarts(targetCartID, sourceCartID string) error {
	return s.mergeCarts(targetCartID, sourceCartID)
}
func (s *CartAPIServer) GetCart(cartID string) (*api.Cart, error) {
	return s.getCart(cartID)
}
func (s *CartAPIServer) AddCartItem(cartID string, cartItem *api.CartItem) error {
	return s.addCartItem(cartID, cartItem)
}

func (s *CartAPIServer) RemoveCartItem(cartID, cartItemID string) error {
	return s.removeCartItem(cartID, cartItemID)
}
func (s *CartAPIServer) UpdateCartItem(cartID string, cartItem *api.CartItem) error {
	return s.updateCartItem(cartID, cartItem)
}
