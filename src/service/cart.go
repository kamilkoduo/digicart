package service

import (
	"github.com/google/uuid"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service/utils"
)

func createEmptyCart(userID string, guestID string) (error) {
	cart := &api.Cart{
		CartID:  uuid.New().String(),
		UserID:  userID,
		GuestID: guestID,
		Items:   make([]*api.CartItem, 0, defaultCartCapacity),
	}
	cartKey := utils.CartKey(cart)
	cartInfoMap := utils.CartInfoMap(cart)
	cartItemsKey :=utils.CartItemsKey(cart)

	redisClient.Del(cartKey)
	redisClient.HMSet(cartKey, cartInfoMap)

	redisClient.Del(cartItemsKey)

}
