package service

import (
	"github.com/kamilkoduo/digicart/src/api"
	"log"
	"strconv"
)

func getCartItems(cartID string) ([]*api.CartItem, error) {
	itemIDs, err := redisClient.SMembers("cart:" + cartID + ":items").Result()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	items := make([]*api.CartItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		offerID, err := redisClient.Get("cart:" + cartID + ":items:" + itemID + ":offer_id").Result()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		offerPriceStr, err := redisClient.Get("cart:" + cartID + ":items:" + itemID + ":offer_price").Result()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		offerPrice, err := strconv.ParseFloat(offerPriceStr, 64)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		offerTitleMap, err := redisClient.HGetAll("cart:" + cartID + ":items:" + itemID + ":offer_title").Result()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		countStr, err := redisClient.Get("cart:" + cartID + ":items:" + itemID + ":count").Result()
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		cartItem := &api.CartItem{
			CartItemID: itemID,
			OfferID:    offerID,
			OfferPrice: offerPrice,
			OfferTitle: offerTitleMap,
			Count:      count,
		}
		items = append(items, cartItem)
	}
	return items, nil
}

//func createEmptyCart(userID string, guestID string) error {
//	cart := &api.Cart{
//		CartID:  uuid.New().String(),
//		UserID:  userID,
//		GuestID: guestID,
//		Items:   make([]*api.CartItem, 0, defaultCartCapacity),
//	}
//	cartKey := utils.CartKey(cart)
//	cartInfoMap := utils.CartInfoMap(cart)
//	cartItemsKey := utils.CartItemsKey(cart)
//
//	redisClient.Del(cartKey)
//	redisClient.HMSet(cartKey, cartInfoMap)
//
//	redisClient.Del(cartItemsKey)
//	return nil
//}
