package service

import (
	"fmt"
	"github.com/kamilkoduo/digicart/src/api"
	"log"
	"strconv"
)

type CartApiServer struct{}

const CartSetKey = "carts"
const KeySep = ":"
const CartTypeKey = "type"
const MergedSetKey = "merged"

func (s *CartApiServer) GetCart(cartID string, cartType api.CartType) (*api.Cart, error) {
	found, err := cartExists(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if !found {
		err := fmt.Errorf("no such cart: %v", cartID)
		log.Printf(err.Error())
		return nil, err
	}
	cartTypeActual, err := getCartType(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if cartType != cartTypeActual{
		err := fmt.Errorf("cart type: %v does not correspond to actual: %v", cartType, cartTypeActual)
		log.Printf(err.Error())
	}
	mergedCartIDs,err := getMergedCartIDs(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	cart := &api.Cart{
		CartID:        cartID,
		CartType:      cartTypeActual,
		MergedCartIDs: mergedCartIDs,
		Items:         nil,
	}
	//items, err := getCartItems(cartID)
	//for _, x := range items {
	//	println(x)
	//}
	//println(err == nil)
	return cart, err
}
//func (s *CartApiServer) AddCartItem(userID string, guestID string) (error)    { return nil }
//func (s *CartApiServer) UpdateCartItem(userID string, guestID string) (error) { return nil }
//func (s *CartApiServer) RemoveCartItem(userID string, guestID string) (error) { return nil }


//func joinCarts(cartSourceID string, cartTargetID string) (string, error) {
//	return "", nil
//}

func cartExists(cartID string) (bool, error) {
	found, err := redisClient.SIsMember(CartSetKey, cartID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return found, err
}

func getCartType(cartID string) (api.CartType, error) {
	cartTypeStr, err := redisClient.Get(CartSetKey + KeySep + cartID + KeySep + CartTypeKey).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	cartType, err := strconv.Atoi(cartTypeStr)
	if err != nil {
		log.Printf(err.Error())
	}
	return (api.CartType)(cartType), err
}
func getMergedCartIDs(cartID string) ([]string, error) {
	mergedCartsSet, err := redisClient.SMembers(CartSetKey + KeySep + cartID + KeySep + MergedSetKey).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return mergedCartsSet, err
}