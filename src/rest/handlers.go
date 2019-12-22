package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"log"
	"net/http"
)

func AuthenticationRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Request.Header.Get("x-user-id")
		guestID := ctx.Request.Header.Get("x-guest-id")

		if userID != "" {
			ctx.Set("cartID", userID)
			ctx.Set("cartType", api.CartType_Authorized)

			if guestID != "" {
				ctx.Set("mergeCartID", guestID)
			}
		} else if guestID != "" {
			ctx.Set("cartID", guestID)
			ctx.Set("cartType", api.CartType_Guest)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		}
		ctx.Next()
	}

}

func CartAuthorization(initCartIfDoesNotExist bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartID, _ := ctx.Get("cartID")
		cartType, _ := ctx.Get("cartType")
		exists, err := cartApi.CartExists((cartID).(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, withFutureMerge := ctx.Get("mergeCartID")
		if exists {
			cartTypeActual, err := cartApi.GetCartType((cartID).(string))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if cartType != cartTypeActual {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				return
			}
		} else if initCartIfDoesNotExist || withFutureMerge {
			err = cartApi.InitCart((cartID).(string), (cartType).(api.CartType))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.Next()
	}
}
func CartMerge() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartID, _ := ctx.Get("cartID")
		mergeCartID, found := ctx.Get("mergeCartID")
		if found {
			err := cartApi.MergeCarts(cartID.(string), mergeCartID.(string))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.Next()
	}
}

func CartItemPreprocess(withPayload bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartItemID := ctx.Param("id")
		if cartItemID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "cart item id was not supplied"})
		}
		if withPayload {
			var cartItem = &api.CartItem{}
			err := jsonapi.UnmarshalPayload(ctx.Request.Body, cartItem)
			if err != nil {
				log.Printf(err.Error())
				ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
				return
			}
			cartItem.CartItemID = cartItemID
			ctx.Set("cartItem", cartItem)
		} else {
			ctx.Set("cartItemID", cartItemID)
		}
		ctx.Next()
	}
}
