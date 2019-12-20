package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service"
	"log"
	"net/http"
)


var cartApi api.CartApi = service.CartApiServer{}

func Run() {
	router := gin.Default()
	router.Use(AuthenticationRequired())
	//get
	router.GET("/",
		CartAuthorization(false),
		CartMerge(),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			cart, err := cartApi.GetCart(cartID.(string))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.Header("content-type", jsonapi.MediaType)
			//ctx.Header("content-type", "application/vnd.cartApi+json; charset=utf-8")
			ctx.Status(http.StatusOK)

			err = jsonapi.MarshalPayload(ctx.Writer, cart)
			if err != nil {
				log.Printf(err.Error())
			}
		})
	//todo: навести порядок с ошибками
	//post
	router.POST("/:id",
		CartAuthorization(true),
		CartMerge(),
		CartItemPreprocess(true),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			cartItem,_ :=ctx.Get("cartItem")
			err := cartApi.AddCartItem(cartID.(string), cartItem.(*api.CartItem))
			if err != nil {
				log.Printf(err.Error())
			}
			ctx.Status(http.StatusOK)
		})
	//put
	router.PUT("/:id",
		CartAuthorization(false),
		CartMerge(),
		CartItemPreprocess(true),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			cartItem,_ :=ctx.Get("cartItem")
			err := cartApi.UpdateCartItem((cartID).(string), cartItem.(*api.CartItem))
			if err != nil {
				log.Printf(err.Error())
			}
			ctx.Status(http.StatusOK)
		})
	//post
	router.DELETE("/:id",
		CartAuthorization(false),
		CartMerge(),
		CartItemPreprocess(false),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			cartItemID,_:= ctx.Get("cartItemID")
			err := cartApi.RemoveCartItem((cartID).(string), (cartItemID).(string))
			if err != nil {
				log.Printf(err.Error())
			}
			ctx.Status(http.StatusOK)
		})


	log.Fatalf("Gin Router failed: %+v", router.Run(service.AppAddress))
}
