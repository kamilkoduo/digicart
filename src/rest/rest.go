package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/rest/config"
	"log"
	"net/http"
)

var cartAPI api.CartAPI = config.APIServer

func Run() {
	router := gin.Default()
	router.Use(JSONAPIHeader())
	router.Use(AuthenticationRequired())

	v1ApiMy := router.Group(config.PathAPIv1My)
	{
		v1ApiMy.GET("/",
			CartAuthorization(false),
			CartMerge(),
			func(ctx *gin.Context) {
				cartID, _ := ctx.Get(config.KeyCartID)
				cart, err := cartAPI.GetCart(cartID.(string))
				if err != nil {
					cErr := err.(carterrors.CartError)
					ctx.Status(cErr.StatusCode())
					_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
					ctx.Abort()
					return
				}
				ctx.Status(http.StatusOK)
				err = jsonapi.MarshalPayload(ctx.Writer, cart)
				if err != nil {
					log.Fatalf(err.Error())
				}
			})
		items := v1ApiMy.Group(config.PathPostfixItems)
		{
			items.POST("/:"+config.KeyID,
				CartAuthorization(true),
				CartMerge(),
				CartItemPreprocess(true),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get(config.KeyCartID)
					cartItem, _ := ctx.Get(config.KeyCartItem)
					err := cartAPI.AddCartItem(cartID.(string), cartItem.(*api.CartItem))
					if err != nil {
						cErr := err.(carterrors.CartError)
						ctx.Status(cErr.StatusCode())
						_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
						ctx.Abort()
						return
					}
					ctx.Status(http.StatusOK)
				})
			// put
			items.PUT("/:"+config.KeyID,
				CartAuthorization(false),
				CartMerge(),
				CartItemPreprocess(true),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get(config.KeyCartID)
					cartItem, _ := ctx.Get(config.KeyCartItem)
					err := cartAPI.UpdateCartItem((cartID).(string), cartItem.(*api.CartItem))
					if err != nil {
						cErr := err.(carterrors.CartError)
						ctx.Status(cErr.StatusCode())
						_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
						ctx.Abort()
						return
					}
					ctx.Status(http.StatusOK)
				})
			// post
			items.DELETE("/:"+config.KeyID,
				CartAuthorization(false),
				CartMerge(),
				CartItemPreprocess(false),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get(config.KeyCartID)
					cartItemID, _ := ctx.Get(config.KeyCartItemID)
					err := cartAPI.RemoveCartItem((cartID).(string), (cartItemID).(string))
					if err != nil {
						cErr := err.(carterrors.CartError)
						ctx.Status(cErr.StatusCode())
						_ = jsonapi.MarshalErrors(ctx.Writer, []*jsonapi.ErrorObject{cErr.JSONObject()})
						ctx.Abort()
						return
					}
					ctx.Status(http.StatusOK)
				})
		}
	}

	log.Fatalf("Gin Router failed: %+v", router.Run(config.AppAddress))
}
