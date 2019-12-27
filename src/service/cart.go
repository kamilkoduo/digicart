package service

import (
	"fmt"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
)

func (s CartAPIServer) getCart(cartID string) (*api.Cart, error) {
	// cart existence is checked in getCartType method
	cartType, err := s.getCartType(cartID)
	if err != nil {
		return nil, err
	}
	mergedCartIDs := s.cartDBAPI.GetMergedCartIDs(cartID)
	items, err := s.getCartItems(cartID)
	if err != nil {
		return nil, err
	}
	cart := &api.Cart{
		CartID:        cartID,
		CartType:      cartType,
		MergedCartIDs: mergedCartIDs,
		Items:         items,
	}
	return cart, nil
}
func (s CartAPIServer) mergeCarts(targetCartID, sourceCartID string) error {
	// check target cart existence
	exists, err := s.CartExists(targetCartID)
	if err != nil {
		return err
	}
	if ! exists {
		return carterrors.New(carterrors.CartNotFound)
	}
	// source cart existence is checked inside getCart method
	sourceCart, err := s.getCart(sourceCartID)
	if err == nil {
		s.cartDBAPI.RemoveCartCompletely(sourceCartID)
		s.cartDBAPI.AddToMergedCartIDs(targetCartID, sourceCart.MergedCartIDs...)
		for _, item := range sourceCart.Items {
			err = s.addCartItem(targetCartID, item)
			if err != nil {
				return err
			}
		}
	} else {
		if !err.(carterrors.CartError).IsType(carterrors.CartNotFound) {
			return err
		}
	}
	s.cartDBAPI.AddToMergedCartIDs(targetCartID, sourceCartID)
	return nil
}
func (s CartAPIServer) cartExists(cartID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	fmt.Printf("%v %T", s.cartDBAPI, s.cartDBAPI)
	found := s.cartDBAPI.CartIDIsPresent(cartID)
	return found, nil
}
func (s CartAPIServer) initCart(cartID string, cartType api.CartType) error {
	exists, err := s.cartExists(cartID)
	if err != nil {
		return err
	}
	if exists {
		return carterrors.New(carterrors.CartAlreadyExists, cartID)
	}
	s.cartDBAPI.AddCartID(cartID)
	s.cartDBAPI.SetCartType(cartID, uint8(cartType))
	return nil

}
func (s CartAPIServer) getCartType(cartID string) (api.CartType, error) {
	found, err := s.cartExists(cartID)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, carterrors.New(carterrors.CartNotFound)
	}
	cartType := s.cartDBAPI.GetCartType(cartID)
	return api.CartType(cartType), nil
}
