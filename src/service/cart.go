package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
)

func getCart(cartID string) (*api.Cart, error) {
	found, err := cartExists(cartID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, carterrors.New(carterrors.CartNotFound, cartID)
	}
	cartType, err := getCartType(cartID)
	if err != nil {
		return nil, err
	}
	mergedCartIDs := cartDBAPI.GetMergedCartIDs(cartID)
	items, err := getCartItems(cartID)
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

func mergeCarts(targetCartID, sourceCartID string) error {
	// check cart existence
	foundS, err := cartExists(sourceCartID)
	if err != nil {
		return err
	}
	if foundS {
		sourceCart, err := getCart(sourceCartID)
		if err != nil {
			return err
		}
		cartDBAPI.RemoveCartCompletely(sourceCartID)
		cartDBAPI.AddToMergedCartIDs(targetCartID, sourceCart.MergedCartIDs...)
		for _, item := range sourceCart.Items {
			err = addCartItem(targetCartID, item)
			if err != nil {
				return err
			}
		}
	}
	cartDBAPI.AddToMergedCartIDs(targetCartID, sourceCartID)
	return nil
}

func cartExists(cartID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	found := cartDBAPI.CartIDIsPresent(cartID)
	return found, nil
}
func initCart(cartID string, cartType api.CartType) error {
	exists, err := cartExists(cartID)
	if err != nil {
		return err
	}
	if exists {
		return carterrors.New(carterrors.CartAlreadyExists, cartID)
	}
	cartDBAPI.AddCartID(cartID)
	cartDBAPI.SetCartType(cartID, uint8(cartType))
	return nil

}



func getCartType(cartID string) (api.CartType, error) {
	found, err := cartExists(cartID)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, carterrors.New(carterrors.CartNotFound)
	}
	cartType := cartDBAPI.GetCartType(cartID)
	return api.CartType(cartType), nil
}
