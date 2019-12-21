package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
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
	addCartItemOfferTitle(cartID, cartItem.CartItemID, ((interface{})(cartItem.OfferTitle)).(map[string]interface{}))
	return nil
}

//func updateCartItemOnlyCount (cartID string, cartItem *api.CartItem) error {
//	cartItemOld, err := getCartItem(cartID, cartItem.CartItemID)
//	if err != nil {
//		log.Printf(err.Error())
//	}
//	if cartItemOld.OfferID != off
//}
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
func removeCartItem(cartID string, cartItemID string) error {
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
func cartItemExists(cartID string, cartItemID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	if !validID(cartItemID) {
		return false, carterrors.New(carterrors.InvalidCartItemID, cartItemID)
	}
	found, err := redisClient.SIsMember(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return found, nil
}

func getCartItemInfo(cartID string, cartItemID string) (map[string]string, error) {
	info, err := redisClient.HGetAll(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return info, err
}
func getCartItemOfferTitle(cartID string, cartItemID string) (map[string]interface{}, error) {
	title, err := redisClient.HGetAll(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	//todo: map-map
	titleSI := make(map[string]interface{})
	for k, v := range title {
		titleSI[k] = v
	}
	return titleSI, err
}
func getCartItem(cartID string, cartItemID string) (*api.CartItem, error) {
	infoMap, err := getCartItemInfo(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	titleMap, err := getCartItemOfferTitle(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	offerID := infoMap[keys.ItemOfferIDMapKey()]
	offerPrice, err := strconv.ParseFloat(infoMap[keys.ItemOfferPriceMapKey()], 64)
	if err != nil {
		log.Printf(err.Error())
	}
	count, err := strconv.ParseUint(infoMap[keys.ItemCountMapKey()], 10, 0)
	if err != nil {
		log.Printf(err.Error())
	}
	cartItem := &api.CartItem{
		CartItemID: cartItemID,
		OfferID:    offerID,
		OfferPrice: offerPrice,
		OfferTitle: titleMap,
		Count:      uint(count),
	}
	return cartItem, err
}
func getCartItems(cartID string) ([]*api.CartItem, error) {
	itemIDs, err := redisClient.SMembers(keys.CartItemSet(cartID)).Result()
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
func addCartItemID(cartID string, cartItemID string) {
	_, err := redisClient.SAdd(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return
}
func removeCartItemID(cartID string, cartItemID string) {
	_, err := redisClient.SRem(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func addCartItemInfo(cartID string, cartItemID string, cartItemInfo map[string]interface{}) {
	_, err := redisClient.HMSet(keys.ItemInfoMap(cartID, cartItemID), cartItemInfo).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func removeCartItemInfo(cartID string, cartItemID string) {
	_, err := redisClient.Del(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func addCartItemOfferTitle(cartID string, cartItemID string, cartItemOfferTitle map[string]interface{}) {
	_, err := redisClient.HMSet(keys.ItemOfferTitleMap(cartID, cartItemID), cartItemOfferTitle).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func removeCartItemOfferTitle(cartID string, cartItemID string) {
	_, err := redisClient.Del(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
