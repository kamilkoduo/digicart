package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service/db/service/keys"
	"log"
	"strconv"
)

func (s CartAPIServer) cartItemExists(cartID, cartItemID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	if !validID(cartItemID) {
		return false, carterrors.New(carterrors.InvalidCartItemID, cartItemID)
	}
	cartFound := s.cartDBAPI.CartIDIsPresent(cartID)
	if !cartFound {
		return false, carterrors.New(carterrors.CartNotFound)
	}
	found := s.cartDBAPI.CartItemIDIsPresent(cartID, cartItemID)
	return found, nil
}
func (s CartAPIServer) getCartItem(cartID, cartItemID string) (*api.CartItem, error) {
	found, err := s.cartItemExists(cartID, cartItemID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, carterrors.New(carterrors.CartItemNotFound, cartID, cartItemID)
	}
	infoMap := s.cartDBAPI.GetCartItemInfo(cartID, cartItemID)
	titleMap := s.cartDBAPI.GetCartItemOfferTitle(cartID, cartItemID)
	offerID := infoMap[keys.ItemOfferIDMapKey()]
	offerPrice, err := strconv.ParseFloat(infoMap[keys.ItemOfferPriceMapKey()], 64)
	if err != nil {
		log.Fatalf(err.Error())
	}
	count, err := strconv.ParseUint(infoMap[keys.ItemCountMapKey()], 10, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	cartItem := &api.CartItem{
		CartItemID: cartItemID,
		OfferID:    offerID,
		OfferPrice: offerPrice,
		OfferTitle: titleMap,
		Count:      uint(count),
	}
	return cartItem, nil
}
func (s CartAPIServer) getCartItems(cartID string) ([]*api.CartItem, error) {
	itemIDs := s.cartDBAPI.GetCartItemIDs(cartID)
	items := make([]*api.CartItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		cartItem, err := s.getCartItem(cartID, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, cartItem)
	}
	return items, nil
}
func (s CartAPIServer) addCartItem(cartID string, cartItem *api.CartItem) error {
	//cart existence is checked in cart item exists
	found, err := s.cartItemExists(cartID, cartItem.CartItemID)
	if err != nil {
		return err
	}
	if found {
		return carterrors.New(carterrors.CartItemAlreadyExists, cartID, cartItem.CartItemID)
	}
	s.cartDBAPI.AddCartItemID(cartID, cartItem.CartItemID)
	infoMap := map[string]interface{}{
		keys.ItemOfferIDMapKey():    cartItem.OfferID,
		keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
		keys.ItemCountMapKey():      cartItem.Count,
	}
	s.cartDBAPI.AddCartItemInfo(cartID, cartItem.CartItemID, infoMap)
	s.cartDBAPI.AddCartItemOfferTitle(cartID, cartItem.CartItemID, (interface{}(cartItem.OfferTitle)).(map[string]interface{}))
	return nil
}
func (s CartAPIServer) updateCartItem(cartID string, cartItem *api.CartItem) error {
	//existence is checked in remove cart item
	err := s.removeCartItem(cartID, cartItem.CartItemID)
	if err != nil {
		return err
	}
	err = s.addCartItem(cartID, cartItem)
	if err != nil {
		return err
	}
	return nil
}
func (s CartAPIServer) removeCartItem(cartID, cartItemID string) error {
	found, err := s.cartItemExists(cartID, cartItemID)
	if err != nil {
		return err
	}
	if !found {
		return carterrors.New(carterrors.CartItemNotFound, cartID, cartItemID)
	}
	s.cartDBAPI.RemoveCartItemID(cartID, cartItemID)
	s.cartDBAPI.RemoveCartItemInfo(cartID, cartItemID)
	s.cartDBAPI.RemoveCartItemOfferTitle(cartID, cartItemID)
	return nil
}
