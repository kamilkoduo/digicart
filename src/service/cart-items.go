package service

import (
	"fmt"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service/keys"
	"log"
	"strconv"
)

func addCartItem(cartID string, cartItem *api.CartItem) error {
	found, err := cartItemExists(cartID, cartItem.CartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	if found {
		err := fmt.Errorf("cart item already exists : `%v`-`%v`", cartID, cartItem.CartItemID)
		log.Printf(err.Error())
		return err
	}
	err = addCartItemID(cartID, cartItem.CartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	infoMap := map[string]interface{}{
		keys.ItemOfferIDMapKey():    cartItem.OfferID,
		keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
		keys.ItemCountMapKey():      cartItem.Count,
	}
	err = addCartItemInfo(cartID, cartItem.CartItemID, infoMap)
	if err != nil {
		log.Printf(err.Error())
	}

	err = addCartItemOfferTitle(cartID, cartItem.CartItemID, ((interface{})(cartItem.OfferTitle)).(map[string]interface{}))
	if err != nil {
		log.Printf(err.Error())
	}
	return err

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
		log.Printf(err.Error())
	}
	if !found {
		err := fmt.Errorf("no such cart item: `%v`-`%v`", cartID, cartItem.CartItemID)
		log.Printf(err.Error())
		return err
	}

	err = removeCartItem(cartID, cartItem.CartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	err = addCartItem(cartID,cartItem)
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func removeCartItem(cartID string, cartItemID string) error {
	found, err := cartItemExists(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	if !found {
		err := fmt.Errorf("no such cart item : `%v`-`%v`", cartID, cartItemID)
		log.Printf(err.Error())
		return err
	}
	err = removeCartItemID(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	err = removeCartItemInfo(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	err = removeCartItemOfferTitle(cartID, cartItemID)
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func cartItemExists(cartID string, cartItemID string) (bool, error) {
	found, err := redisClient.SIsMember(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return found, err
}

func getCartItemInfo(cartID string, cartItemID string) (map[string]string, error) {
	info, err := redisClient.HGetAll(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return info, err
}
func getCartItemOfferTitle(cartID string, cartItemID string) (map[string]string, error) {
	titles, err := redisClient.HGetAll(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return titles, err
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
	count, err := strconv.ParseUint(infoMap[keys.ItemCountMapKey()], 10 ,0)
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
		log.Printf(err.Error())
		return nil, err
	}
	items := make([]*api.CartItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		cartItem, err := getCartItem(cartID, itemID)
		if err != nil {
			log.Printf(err.Error())
		}
		items = append(items, cartItem)
	}
	return items, err
}
func addCartItemID(cartID string, cartItemID string) error {
	_, err := redisClient.SAdd(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func removeCartItemID(cartID string, cartItemID string) error {
	_, err := redisClient.SRem(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func addCartItemInfo(cartID string, cartItemID string, cartItemInfo map[string]interface{}) error {
	_, err := redisClient.HMSet(keys.ItemInfoMap(cartID, cartItemID), cartItemInfo).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func removeCartItemInfo(cartID string, cartItemID string) error {
	_, err := redisClient.Del(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}

func addCartItemOfferTitle(cartID string, cartItemID string, cartItemOfferTitle map[string]interface{}) error {
	_, err := redisClient.HMSet(keys.ItemOfferTitleMap(cartID, cartItemID), cartItemOfferTitle).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
func removeCartItemOfferTitle(cartID string, cartItemID string) error {
	_, err := redisClient.Del(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
