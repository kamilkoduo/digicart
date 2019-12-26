package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service/db/service/keys"
	"log"
	"strconv"
)

func addCartItem(cartID string, cartItem *api.CartItem) error {
	found, err := cartItemExists(cartID, cartItem.CartItemID)
	if err != nil {
		return err
	}
	if found {
		return carterrors.New(carterrors.CartItemAlreadyExists, cartID, cartItem.CartItemID)
	}
	cartDBAPI.AddCartItemID(cartID, cartItem.CartItemID)
	infoMap := map[string]interface{}{
		keys.ItemOfferIDMapKey():    cartItem.OfferID,
		keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
		keys.ItemCountMapKey():      cartItem.Count,
	}
	cartDBAPI.AddCartItemInfo(cartID, cartItem.CartItemID, infoMap)
	cartDBAPI.AddCartItemOfferTitle(cartID, cartItem.CartItemID, (interface{}(cartItem.OfferTitle)).(map[string]interface{}))
	return nil
}

func updateCartItem(cartID string, cartItem *api.CartItem) error {
	found, err := cartItemExists(cartID, cartItem.CartItemID)
	if err != nil {
		return err
	}
	if !found {
		return carterrors.New(carterrors.CartItemNotFound, cartID, cartItem.CartItemID)
	}
	err = removeCartItem(cartID, cartItem.CartItemID)
	if err != nil {
		return err
	}
	err = addCartItem(cartID, cartItem)
	if err != nil {
		return err
	}
	return nil
}
func removeCartItem(cartID, cartItemID string) error {
	found, err := cartItemExists(cartID, cartItemID)
	if err != nil {
		return err
	}
	if !found {
		return carterrors.New(carterrors.CartItemNotFound, cartID, cartItemID)
	}
	cartDBAPI.RemoveCartItemID(cartID, cartItemID)
	cartDBAPI.RemoveCartItemInfo(cartID, cartItemID)
	cartDBAPI.RemoveCartItemOfferTitle(cartID, cartItemID)
	return nil
}
func cartItemExists(cartID, cartItemID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	if !validID(cartItemID) {
		return false, carterrors.New(carterrors.InvalidCartItemID, cartItemID)
	}
	found := cartDBAPI.CartItemIDIsPresent(cartID, cartItemID)
	return found, nil
}

func getCartItem(cartID, cartItemID string) (*api.CartItem, error) {
	found, err := cartItemExists(cartID, cartItemID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, carterrors.New(carterrors.CartItemNotFound, cartID, cartItemID)
	}
	infoMap := cartDBAPI.GetCartItemInfo(cartID, cartItemID)
	titleMap := cartDBAPI.GetCartItemOfferTitle(cartID, cartItemID)
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
func getCartItems(cartID string) ([]*api.CartItem, error) {
	itemIDs := cartDBAPI.GetCartItemIDs(cartID)
	items := make([]*api.CartItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		cartItem, err := getCartItem(cartID, itemID)
		if err != nil {
			return nil, err
		}
		items = append(items, cartItem)
	}
	return items, nil
}
