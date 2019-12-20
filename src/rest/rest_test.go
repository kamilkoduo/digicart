package rest_test

import (
	"github.com/google/uuid"
	"github.com/kamilkoduo/digicart/src/rest"
	"github.com/kamilkoduo/digicart/src/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"time"
)

var _ = Describe("Rest", func() {
	const userHeader = "x-user-id"
	data := struct {
		user1 string
	}{
		uuid.New().String(),
	}
	var url struct {
		app      string
		cart     string
		cartItem func(id string) string
	}
	BeforeSuite(func() {
		//run server
		func() {
			go rest.Run()
			time.Sleep(1.5e9)
		}()
		//init urls
		func() {
			url.app = "http://"+service.AppAddress
			url.cart = url.app + "/api/v1/cart-api/my"
			url.cartItem = func(id string) string {
				return url.cart + "/items/:" + id
			}
		}()
		//
	})

	Describe("Interaction with Cart API", func() {
		Describe("Authorization", func() {
			Context("No auth headers provided", func() {
				It("should reply", func() {
					resp, err := http.Get(url.cart)
					log.Print(url)
					if err!=nil{
						log.Printf("error %v",err.Error())
					}
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})
			PContext("User auth header provided", func() {
				It("should reply", func() {
					client := &http.Client{}
					req, err := http.NewRequest("GET", url.cart, nil)
					Expect(err).ShouldNot(HaveOccurred())
					req.Header.Set(userHeader, data.user1)
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Status).To(Equal(http.StatusOK))
				})
			})
		})
	})
})
