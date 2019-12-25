package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/rest/config"
)

func JSONAPIHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header(config.HeaderContentType, jsonapi.MediaType)
	}
}
func AuthenticationRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Request.Header.Get(config.HeaderUserID)
		guestID := ctx.Request.Header.Get(config.HeaderGuestID)

		if userID != "" {
			ctx.Set(config.KeyCartID, userID)
			ctx.Set(config.KeyCartType, api.CartTypeAuthorized)

			if guestID != "" {
				ctx.Set(config.KeyMergeCartID, guestID)
			}
		} else if guestID != "" {
			ctx.Set(config.KeyCartID, guestID)
			ctx.Set(config.KeyCartType, api.CartTypeGuest)
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
		cartID, _ := ctx.Get(config.KeyCartID)
		cartType, _ := ctx.Get(config.KeyCartType)
		exists, err := cartAPI.CartExists((cartID).(string))
		if err != nil {
			cErr := err.(carterrors.CartError)
			ctx.Status(cErr.StatusCode())
			_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
			ctx.Abort()
			return
		}
		_, withFutureMerge := ctx.Get(config.KeyMergeCartID)
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
		cartID, _ := ctx.Get(config.KeyCartID)
		mergeCartID, found := ctx.Get(config.KeyMergeCartID)
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
		cartItemID := ctx.Param(config.KeyID)
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
			ctx.Set(config.KeyCartItem, cartItem)
		} else {
			ctx.Set(config.KeyCartItemID, cartItemID)
		}
		ctx.Next()
	}
}
