package service

import (
	"github.com/kamilkoduo/digicart/src/service/db/service/keys"
	"log"
	"strconv"
)

type CartDBRedisServer struct{}

func (s CartDBRedisServer) GetMergedCartIDs(cartID string) []string {
	mergedCartsSet, err := redisClient.SMembers(keys.CartMergedSet(cartID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return mergedCartsSet
}
func (s CartDBRedisServer) AddCartItemID(cartID, cartItemID string) {
	_, err := redisClient.SAdd(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func (s CartDBRedisServer) RemoveCartItemID(cartID, cartItemID string) {
	_, err := redisClient.SRem(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func (s CartDBRedisServer) AddCartItemInfo(cartID, cartItemID string, cartItemInfo map[string]interface{}) {
	_, err := redisClient.HMSet(keys.ItemInfoMap(cartID, cartItemID), cartItemInfo).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func (s CartDBRedisServer) RemoveCartItemInfo(cartID, cartItemID string) {
	_, err := redisClient.Del(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s CartDBRedisServer) AddCartItemOfferTitle(cartID, cartItemID string, cartItemOfferTitle map[string]interface{}) {
	_, err := redisClient.HMSet(keys.ItemOfferTitleMap(cartID, cartItemID), cartItemOfferTitle).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func (s CartDBRedisServer) RemoveCartItemOfferTitle(cartID, cartItemID string) {
	_, err := redisClient.Del(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func (s CartDBRedisServer) CartIDIsPresent(cartID string) bool{
	found, err := redisClient.SIsMember(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return found
}
func (s CartDBRedisServer) AddCartID(cartID string)  {
	_, err := redisClient.SAdd(keys.CartSet(), cartID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s CartDBRedisServer)SetCartType(cartID string, cartType uint8){
	_, err := redisClient.Set(keys.CartTypeKey(cartID), cartType, -1).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s CartDBRedisServer) AddToMergedCartIDs(cartID string, mergedID ...string) {
	if len(mergedID) > 0 {
		_, err := redisClient.SAdd(keys.CartMergedSet(cartID), mergedID).Result()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
func (s CartDBRedisServer) GetCartType(cartID string) (uint8) {
	cartTypeStr, err := redisClient.Get(keys.CartTypeKey(cartID)).Result()
	if err !=nil {
		log.Fatalf(err.Error())
	}
	cartType, err := strconv.ParseUint(cartTypeStr, 10, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return uint8(cartType)
}
func (s CartDBRedisServer)	CartItemIDIsPresent(cartID, cartItemID string) bool{
	found, err := redisClient.SIsMember(keys.CartItemSet(cartID), cartItemID).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return found
}

func (s CartDBRedisServer) GetCartItemInfo(cartID, cartItemID string) map[string]string {
	info, err := redisClient.HGetAll(keys.ItemInfoMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return info
}
func (s CartDBRedisServer) GetCartItemOfferTitle(cartID, cartItemID string) map[string]interface{} {
	title, err := redisClient.HGetAll(keys.ItemOfferTitleMap(cartID, cartItemID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	titleSI := make(map[string]interface{})
	for k, v := range title {
		titleSI[k] = v
	}
	return titleSI
}

func (s CartDBRedisServer) GetCartItemIDs(cartID string) []string {
	itemIDs, err := redisClient.SMembers(keys.CartItemSet(cartID)).Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return itemIDs
}

func (s CartDBRedisServer) RemoveCartCompletely(cartID string) {
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