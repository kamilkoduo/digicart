package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/service"
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
		} else if guestID != "" {
			ctx.Set("cartID", guestID)
			ctx.Set("cartType", api.CartType_Guest)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		}
		ctx.Next()
	}
}
func AutomaticCartCreation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartID, _ := ctx.Get("cartID")
		cartType, _ := ctx.Get("cartType")
		exists, err := cartApi.CartExists((cartID).(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			err = cartApi.InitCart((cartID).(string), (cartType).(api.CartType))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.Next()
	}
}

var cartApi api.CartApi = service.CartApiServer{}

func Run() {
	router := gin.Default()
	router.Use(AuthenticationRequired())
	//get
	router.GET("/",
		AutomaticCartCreation(),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			cartType, _ := ctx.Get("cartType")
			cart, err := cartApi.GetCart((cartID).(string), (cartType).(api.CartType))
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
		AutomaticCartCreation(),
		func(ctx *gin.Context) {
			cartID, _ := ctx.Get("cartID")
			/*
				var cartIt = &api.CartItem{}

				cartIt.OfferID="1"
				cartIt.CartItemID="s"
				cartIt.OfferTitle= map[string]string{"en":"x","ru":"f"}
				cartIt.Count=1
				cartIt.OfferPrice=400

				_ = jsonapi.MarshalPayload(ctx.Writer, cartIt)
			*/

			var cartItem = &api.CartItem{}
			err := jsonapi.UnmarshalPayload(ctx.Request.Body, cartItem)
			if err != nil {
				log.Printf(err.Error())
				ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
				return
			}
			err = cartApi.AddCartItem((cartID).(string), cartItem)
			if err != nil {
				log.Printf(err.Error())
			}
			ctx.Status(http.StatusOK)

		})

	log.Fatalf("Gin Router failed: %+v", router.Run())
}
