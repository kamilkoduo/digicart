package service

import (
	"fmt"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service/keys"
	"log"
	"strconv"
)

func getCart(cartID string, cartType api.CartType) (*api.Cart, error) {
	found, err := cartExists(cartID)
	if err != nil {
		log.Fatalf("cartExists(): %v", err.Error())
	}
	if !found {
		err := fmt.Errorf("no such cart: %v", cartID)
		log.Printf(err.Error())
		return nil, err
	}
	log.Printf("cart found: %v", cartID)

	cartTypeActual, err := getCartType(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if cartType != cartTypeActual {
		err := fmt.Errorf("cart type: %v does not correspond to actual: %v", cartType, cartTypeActual)
		log.Printf(err.Error())
	}
	mergedCartIDs, err := getMergedCartIDs(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	items, err := getCartItems(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	cart := &api.Cart{
		CartID:        cartID,
		CartType:      cartTypeActual,
		MergedCartIDs: mergedCartIDs,
		Items:         items,
	}
	return cart, err
}

/*
func MergeCarts(targetCartID string, sourceCartID string) error {
	// check carts existence
	foundT, err := cartExists(targetCartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if !foundT {
		err := fmt.Errorf("no such cart (target): %v", targetCartID)
		log.Printf(err.Error())
		return err
	}
	foundS, err := cartExists(sourceCartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if !foundS {
		err := fmt.Errorf("no such cart (source): %v", sourceCartID)
		log.Printf(err.Error())
		return err
	}
	// check user-guest rule
	typeT,err := getCartType(targetCartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if typeT != api.CartType_Authorized {
		err := fmt.Errorf("cart (target) `%v` is not authorized: `%v`", targetCartID, typeT)
		log.Printf(err.Error())
		return err
	}
	typeS,err:= getCartType(sourceCartID)
	if err != nil {
		log.Printf(err.Error())
	}
	if typeS != api.CartType_Guest {
		err := fmt.Errorf("cart (source) `%v` is not guest: `%v`", sourceCartID, typeS)
		log.Printf(err.Error())
		return err
	}


}
*/
//func joinCarts(cartSourceID string, cartTargetID string) (string, error) {
//	return "", nil
//}

func cartExists(cartID string) (bool, error) {
	found, err := redisClient.SIsMember(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return found, err
}
func initCart(cartID string, cartType api.CartType) error {
	_, err := redisClient.SAdd(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	_, err = redisClient.Set(keys.CartTypeKey(cartID), (uint8)(cartType), -1).Result()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}

func getCartType(cartID string) (api.CartType, error) {
	cartTypeStr, err := redisClient.Get(keys.CartTypeKey(cartID)).Result()
	if err != nil {
		log.Printf("getCartType: %v", err.Error())
	}
	fmt.Printf("%T\n",cartTypeStr)
	cartType, err := strconv.ParseUint(cartTypeStr,10,0)
	if err != nil {
		log.Printf(err.Error())
	}
	return (api.CartType)(cartType), err
}

func getMergedCartIDs(cartID string) ([]string, error) {
	mergedCartsSet, err := redisClient.SMembers(keys.CartMergedSet(cartID)).Result()
	if err != nil {
		log.Printf("getMergedCartIDs: %v",err.Error())
	}
	return mergedCartsSet, err
}
