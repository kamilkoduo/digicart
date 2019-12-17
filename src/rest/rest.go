package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"log"
	"net/http"
	"github.com/kamilkoduo/digicart/src/api"
)

//todo: нормально ли использовать строки для ID?

func Run() {
	router := gin.Default()
	//get
	router.GET("/:id", func(ctx *gin.Context) {
		ci := &api.CartItem{
			CartItemID: uuid.New().String(),
			OfferID:    "abcd",
			OfferPrice: 0,
			OfferTitle: map[string]string{"en": "abcd", "ru": "АБСД"},
			Count:      1,
		}

		cart:=&api.Cart{
			CartID:  uuid.New().String(),
			Items:  []*api.CartItem{ci},
		}

		ctx.Header("content-type", "application/vnd.api+json; charset=utf-8")
		ctx.Status(http.StatusOK)

		err := jsonapi.MarshalPayload(ctx.Writer, cart)
		if err != nil {
			log.Fatal(err)
		}
	})
	//post
	router.POST("/:id", func(ctx *gin.Context) {
		//jsonapi.UnmarshalManyPayload()
	})

	log.Fatalf("Gin Router failed: %+v", router.Run())
}
