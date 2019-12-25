package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service"
	"log"
	"net/http"
)

var cartAPI api.CartAPI = service.CartAPIServer{}

func Run() {
	router := gin.Default()
	router.Use(JSONAPIHeader())
	router.Use(AuthenticationRequired())

	v1ApiMy := router.Group("/api/v1/cart-api/my")
	{
		v1ApiMy.GET("/",
			CartAuthorization(false),
			CartMerge(),
			func(ctx *gin.Context) {
				cartID, _ := ctx.Get("cartID")
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
		items := v1ApiMy.Group("/items")
		{
			items.POST("/:id",
				CartAuthorization(true),
				CartMerge(),
				CartItemPreprocess(true),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get("cartID")
					cartItem, _ := ctx.Get("cartItem")
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
			items.PUT("/:id",
				CartAuthorization(false),
				CartMerge(),
				CartItemPreprocess(true),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get("cartID")
					cartItem, _ := ctx.Get("cartItem")
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
			items.DELETE("/:id",
				CartAuthorization(false),
				CartMerge(),
				CartItemPreprocess(false),
				func(ctx *gin.Context) {
					cartID, _ := ctx.Get("cartID")
					cartItemID, _ := ctx.Get("cartItemID")
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

	log.Fatalf("Gin Router failed: %+v", router.Run(service.AppAddress))
}
