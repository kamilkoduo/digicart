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
					if err.(carterrors.CartError).IsType(carterrors.CartNotFound) {
						ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
					} else {
						ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					}
					return
				}
				ctx.Header("content-type", jsonapi.MediaType)
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
						if err.(carterrors.CartError).IsType(carterrors.CartItemAlreadyExists) {
							ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
						} else {
							ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						}
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
						if err.(carterrors.CartError).IsType(carterrors.CartItemNotFound) {
							ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
						} else {
							ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						}
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
						if err.(carterrors.CartError).IsType(carterrors.CartItemNotFound) {
							ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
						} else {
							ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						}
						return
					}
					ctx.Status(http.StatusOK)
				})
		}
	}

	log.Fatalf("Gin Router failed: %+v", router.Run(service.AppAddress))
}
