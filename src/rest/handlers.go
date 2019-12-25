package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
)

func JSONAPIHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("content-type", jsonapi.MediaType)
	}
}
func AuthenticationRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Request.Header.Get("x-user-id")
		guestID := ctx.Request.Header.Get("x-guest-id")

		if userID != "" {
			ctx.Set("cartID", userID)
			ctx.Set("cartType", api.CartTypeAuthorized)

			if guestID != "" {
				ctx.Set("mergeCartID", guestID)
			}
		} else if guestID != "" {
			ctx.Set("cartID", guestID)
			ctx.Set("cartType", api.CartTypeGuest)
		} else {
			cErr := carterrors.New(carterrors.Unauthenticated, "Authentication Required")
			ctx.Status(cErr.StatusCode())
			_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CartAuthorization(initCartIfDoesNotExist bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartID, _ := ctx.Get("cartID")
		cartType, _ := ctx.Get("cartType")
		exists, err := cartAPI.CartExists((cartID).(string))
		if err != nil {
			cErr := err.(carterrors.CartError)
			ctx.Status(cErr.StatusCode())
			_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
			ctx.Abort()
			return
		}
		_, withFutureMerge := ctx.Get("mergeCartID")
		if exists {
			cartTypeActual, err := cartAPI.GetCartType((cartID).(string))
			if err != nil {
				cErr := err.(carterrors.CartError)
				ctx.Status(cErr.StatusCode())
				_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
				ctx.Abort()
				return
			}
			if cartType != cartTypeActual {
				cErr := carterrors.New(carterrors.AccessDenied)
				ctx.Status(cErr.StatusCode())
				_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
				ctx.Abort()
				return
			}
		} else if initCartIfDoesNotExist || withFutureMerge {
			err = cartAPI.InitCart((cartID).(string), (cartType).(api.CartType))
			if err != nil {
				cErr := err.(carterrors.CartError)
				ctx.Status(cErr.StatusCode())
				_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
				ctx.Abort()
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
			err := cartAPI.MergeCarts(cartID.(string), mergeCartID.(string))
			if err != nil {
				cErr := err.(carterrors.CartError)
				ctx.Status(cErr.StatusCode())
				_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
				ctx.Abort()
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
			cErr := carterrors.New(carterrors.InvalidCartItemID)
			ctx.Status(cErr.StatusCode())
			_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
			ctx.Abort()
			return
		}
		if withPayload {
			var cartItem = &api.CartItem{}
			err := jsonapi.UnmarshalPayload(ctx.Request.Body, cartItem)
			if err != nil {
				cErr := err.(carterrors.CartError)
				ctx.Status(cErr.StatusCode())
				_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
				ctx.Abort()
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
