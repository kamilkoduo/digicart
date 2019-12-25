package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service/config"
	"github.com/kamilkoduo/digicart/src/service/keys"
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
	addCartItemID(cartID, cartItem.CartItemID)
	infoMap := map[string]interface{}{
		keys.ItemOfferIDMapKey():    cartItem.OfferID,
		keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
		keys.ItemCountMapKey():      cartItem.Count,
	}
	addCartItemInfo(cartID, cartItem.CartItemID, infoMap)
	addCartItemOfferTitle(cartID, cartItem.CartItemID, (interface{}(cartItem.OfferTitle)).(map[string]interface{}))
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
	removeCartItemID(cartID, cartItemID)
	removeCartItemInfo(cartID, cartItemID)
	removeCartItemOfferTitle(cartID, cartItemID)
	return nil
}
func cartItemExists(cartID, cartItemID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	if !validID(cartItemID) {
		return false, carterrors.New(carterrors.InvalidCartItemID, cartItemID)
	}
	found, err := config.RedisClient.SIsMember(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return found, nil
}

func getCartItemInfo(cartID, cartItemID string) map[string]string {
	info, err := config.RedisClient.HGetAll(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return info
}
func getCartItemOfferTitle(cartID, cartItemID string) map[string]interface{} {
	title, err := config.RedisClient.HGetAll(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	titleSI := make(map[string]interface{})
	for k, v := range title {
		titleSI[k] = v
	}
	return titleSI
}
func getCartItem(cartID, cartItemID string) (*api.CartItem, error) {
	found, err := cartItemExists(cartID, cartItemID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, carterrors.New(carterrors.CartItemNotFound, cartID, cartItemID)
	}
	infoMap := getCartItemInfo(cartID, cartItemID)
	titleMap := getCartItemOfferTitle(cartID, cartItemID)
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
	itemIDs, err := config.RedisClient.SMembers(keys.CartItemSet(cartID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
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
func addCartItemID(cartID, cartItemID string) {
	_, err := config.RedisClient.SAdd(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func removeCartItemID(cartID, cartItemID string) {
	_, err := config.RedisClient.SRem(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func addCartItemInfo(cartID, cartItemID string, cartItemInfo map[string]interface{}) {
	_, err := config.RedisClient.HMSet(keys.ItemInfoMap(cartID, cartItemID), cartItemInfo).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func removeCartItemInfo(cartID, cartItemID string) {
	_, err := config.RedisClient.Del(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func addCartItemOfferTitle(cartID, cartItemID string, cartItemOfferTitle map[string]interface{}) {
	_, err := config.RedisClient.HMSet(keys.ItemOfferTitleMap(cartID, cartItemID), cartItemOfferTitle).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func removeCartItemOfferTitle(cartID, cartItemID string) {
	_, err := config.RedisClient.Del(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
