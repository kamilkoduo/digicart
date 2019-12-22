package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kamilkoduo/digicart/src/rest"
	"github.com/kamilkoduo/digicart/src/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var _ = Describe("Rest", func() {
	const (
		userHeader  = "x-user-id"
		guestHeader = "x-guest-id"
	)
	dataIDs := struct {
		user1Existed     string
		guest1Existed    string
		user2NotExisted  string
		guest2NotExisted string
		cartItemUser1    string
		cartitemGuest1   string
		cartItemUser2    string
		cartitemGuest2   string
	}{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}
	cartItemUserData1 := []byte(`{"data": {"type": "cart_item","attributes": {"count": 1,"offer_id": "1",
							"offer_price": 400,"offer_title": {"en": "x","ru": "f"}}}}`)
	cartItemUserData2 := []byte(`{"data": {"type": "cart_item","attributes": {"count": 1,"offer_id": "1",
							"offer_price": 200,"offer_title": {"en": "x","ru": "f"}}}}`)
	cartItemGuestData1 := []byte(`{"data": {"type": "cart_item","attributes": {"count": 2,"offer_id": "10",
							"offer_price": 1400,"offer_title": {"en": "y","ru": "h"}}}}`)
	cartItemGuestData2 := []byte(`{"data": {"type": "cart_item","attributes": {"count": 2,"offer_id": "10",
							"offer_price": 700,"offer_title": {"en": "y","ru": "h"}}}}`)

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
			url.app = "http://" + service.AppAddress
			url.cart = url.app + "/api/v1/cart-api/my"
			url.cartItem = func(id string) string {
				return url.cart + "/items/:" + id
			}
		}()
		////populate cart user 1
		//func() {
		//	client := &http.Client{}
		//	req, _ := http.NewRequest("GET", url.cartItem(dataIDs.user1Existed), nil)
		//	req.Header.Set(userHeader, dataIDs.guest1Existed)
		//	resp, err := client.Do(req)
		//	Expect(err).ShouldNot(HaveOccurred())
		//	Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
		//}()
		/*//populate cart item
		func() {
			jsonData := []byte(`{"data": {"type": "cart_item","attributes": {"count": 1,"offer_id": "1",
							"offer_price": 400,"offer_title": {"en": "x","ru": "f"}}}}`)

			var v interface{}
			_ = json.Unmarshal(jsonData, &v)
			cartItemData = v.(map[string]interface{})
		}()*/
	})

	Describe("Interaction with Cart API", func() {
		Describe("Authorized interaction", func() {
			Context("User cart", func() {
				Context("Get uncreated yet cart", func() {
					Context("cart was not created", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							log.Printf("1111111111")
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
				})
				Context("Add a valid cart item", func() {
					Context("cart item does not exist yet", func() {
						It("should reply OK 1", func() {
							client := &http.Client{}
							req, err := http.NewRequest("POST", url.cartItem(dataIDs.cartItemUser1), bytes.NewBuffer(cartItemUserData1))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
						It("should reply OK 2", func() {
							client := &http.Client{}
							req, err := http.NewRequest("POST", url.cartItem(dataIDs.cartItemUser2), bytes.NewBuffer(cartItemUserData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
					Context("cart item already exists", func() {
						It("should block", func() {
							client := &http.Client{}
							req, err := http.NewRequest("POST", url.cartItem(dataIDs.cartItemUser2), bytes.NewBuffer(cartItemUserData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusConflict))
						})
					})
				})
				Context("Remove a cart item", func() {
					Context("cart item exists", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("DELETE", url.cartItem(dataIDs.cartItemUser2), bytes.NewBuffer(cartItemUserData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
					Context("cart item does not exist", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("DELETE", url.cartItem(dataIDs.cartItemUser2), bytes.NewBuffer(cartItemUserData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
				})
				Context("Update a cart item", func() {
					Context("cart item exists", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("PUT", url.cartItem(dataIDs.cartItemUser1), bytes.NewBuffer(cartItemUserData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
					Context("cart item does not exist", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("PUT", url.cartItem(dataIDs.cartItemUser2), bytes.NewBuffer(cartItemUserData1))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
				})
				Context("Get cart", func() {
					Context("cart exists", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
				})
			})
			Context("Guest cart", func() {
				Context("Add a valid cart item", func() {
					Context("cart item does not exist yet", func() {
						It("should reply OK 1", func() {
							client := &http.Client{}
							req, err := http.NewRequest("POST", url.cartItem(dataIDs.cartitemGuest1), bytes.NewBuffer(cartItemGuestData1))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
						It("should reply OK 2", func() {
							client := &http.Client{}
							req, err := http.NewRequest("POST", url.cartItem(dataIDs.cartitemGuest2), bytes.NewBuffer(cartItemGuestData2))
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
				})
				Context("Get cart", func() {
					Context("cart exists", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
				})
			})

		})
		Describe("Authorization issues", func() {
			Context("No auth headers provided", func() {
				It("should block request", func() {
					resp, err := http.Get(url.cart)
					log.Print(url)
					if err != nil {
						log.Printf("error %v", err.Error())
					}
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})
			Context("User auth header provided", func() {
				Context("Cart exists", func() {
					Context("Cart type is `guest`", func() {
						It("should block request", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
						})
					})
					Context("Cart type is `user`", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
				})
				Context("Cart does not exist", func() {
					Context("Cart type is `guest`", func() {
						It("should block request", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.guest2NotExisted)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
					Context("Cart type is `user`", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user2NotExisted)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
				})
			})
			Context("Guest auth header provided", func() {
				Context("Cart exists", func() {
					Context("Cart type is `guest`", func() {
						It("should block request", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
					Context("Cart type is `user`", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
						})
					})
				})
				Context("Cart does not exist", func() {
					Context("Cart type is `guest`", func() {
						It("should block request", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest2NotExisted)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
					Context("Cart type is `user`", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.user2NotExisted)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
				})
			})
		})
		Describe("Authorized merge", func() {
			Context("Merge carts", func() {
				Context("Get cart", func() {
					Context("cart exists", func() {
						It("should reply OK", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
						})
					})
					Context("guest cart does not exist anymore", func() {
						It("should reply NotFound", func() {
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(guestHeader, dataIDs.guest1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
						})
					})
					Context("user cart contains guest ID", func() {
						It("should reply OK", func() {
							// cart get
							client := &http.Client{}
							req, err := http.NewRequest("GET", url.cart, nil)
							Expect(err).ShouldNot(HaveOccurred())
							req.Header.Set(userHeader, dataIDs.user1Existed)
							resp, err := client.Do(req)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))

							//check guest id
							body, err := ioutil.ReadAll(resp.Body)
							Expect(err).ShouldNot(HaveOccurred())
							jsonMap := make(map[string]interface{})
							err = json.Unmarshal(body, &jsonMap)
							Expect(err).ShouldNot(HaveOccurred())
							guestList := jsonMap["data"].(map[string]interface{})["attributes"].(map[string]interface{})["merged_cart_ids"].([]interface{})
							fmt.Print("LIST ", jsonMap["data"].(map[string]interface{})["attributes"].(map[string]interface{})["merged_cart_ids"])
							contains := func(list []interface{}, element string) bool {
								for _, x := range list {
									if strings.Compare(x.(string), element) == 0 {
										return true
									}
								}
								return false
							}(guestList, dataIDs.guest1Existed)
							Expect(contains).To(BeTrue())
						})
					})
				})
			})
		})
	})
})
