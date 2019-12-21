package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service/keys"
	"log"
	"strconv"
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
	mergedCartIDs := getMergedCartIDs(cartID)
	items, err := getCartItems(cartID)
	if err != nil {
		log.Printf(err.Error())
	}
	cart := &api.Cart{
		CartID:        cartID,
		CartType:      cartType,
		MergedCartIDs: mergedCartIDs,
		Items:         items,
	}
	return cart, nil
}

func mergeCarts(targetCartID string, sourceCartID string) error {
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
		removeCart(sourceCartID)
		addToMergedCartIDs(targetCartID, sourceCart.MergedCartIDs...)
		if err != nil {
			log.Printf(err.Error())
		}
		for _, item := range sourceCart.Items {
			err = addCartItem(targetCartID, item)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}
	addToMergedCartIDs(targetCartID, sourceCartID)
	return nil
}

func cartExists(cartID string) (bool, error) {
	if !validID(cartID) {
		return false, carterrors.New(carterrors.InvalidCartID, cartID)
	}
	found, err := redisClient.SIsMember(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
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
	_, err = redisClient.SAdd(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	_, err = redisClient.Set(keys.CartTypeKey(cartID), (uint8)(cartType), -1).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return nil
}

func removeCart(cartID string) {
	_, err := redisClient.SRem(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	keysToRemove, err := redisClient.Keys(keys.CartPrefix(cartID) + "*").Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	_, err = redisClient.Del(keysToRemove...).Result()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getCartType(cartID string) (api.CartType, error) {
	cartTypeStr, err := redisClient.Get(keys.CartTypeKey(cartID)).Result()
	if err == redis.Nil {
		return 0, carterrors.New(carterrors.CartNotFound)
	}
	cartType, err := strconv.ParseUint(cartTypeStr, 10, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return (api.CartType)(cartType), nil
}

func getMergedCartIDs(cartID string) []string {
	mergedCartsSet, err := redisClient.SMembers(keys.CartMergedSet(cartID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return mergedCartsSet
}

func addToMergedCartIDs(cartID string, mergedID ...string) {
	fmt.Print("ATTENTION: adding these merged ids: ", mergedID, "\n")
	if len(mergedID) > 0 {
		_, err := redisClient.SAdd(keys.CartMergedSet(cartID), mergedID).Result()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
