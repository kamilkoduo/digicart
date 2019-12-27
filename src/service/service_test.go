//+build unit

package service_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/kamilkoduo/digicart/src/api"
	"github.com/kamilkoduo/digicart/src/carterrors"
	"github.com/kamilkoduo/digicart/src/service"
	"github.com/kamilkoduo/digicart/src/service/db/service/keys"
	"github.com/kamilkoduo/digicart/src/service/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("Cart Service", func() {
	var (
		// utils
		ctrl *gomock.Controller
		t    GinkgoTestReporter
		s    api.CartAPI
		dbs  *mocks.MockCartDBAPI
	)
	// for cart test
	const ()
	// ended
	BeforeEach(func() {
		ctrl = gomock.NewController(t)
		dbs = mocks.NewMockCartDBAPI(ctrl)
		srv := &service.CartAPIServer{}
		srv.RegisterCartDBAPIServer(dbs)
		s = srv
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Unit Testing", func() {
		Describe("Cart API methods", func() {
			Describe("Get Cart", func() {
				Context("Invalid cart ID", func() {
					It("", func() {
						//test data
						const (
							cartID   = ""
						)
						_, err := s.GetCart(cartID)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.InvalidCartID))
					})
				})
				Context("Cart exists", func() {
					It("returns cart", func() {
						//test data
						const (
							cartID     = "1"
							cartType   = api.CartTypeAuthorized
							cartItemID = "101"
							offerPrice = float64(4000)
							count      = uint(1)
							offerID    = "10101"
						)
						var offerTitle = map[string]interface{}{"ru": "Элемент корзины 101"}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(2)
						dbs.EXPECT().GetCartType(cartID).Return(uint8(cartType)).Times(1)
						dbs.EXPECT().GetMergedCartIDs(cartID).Return([]string{}).Times(1)
						dbs.EXPECT().GetCartItemIDs(cartID).Return([]string{cartItemID}).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItemID).Return(true).Times(1)
						dbs.EXPECT().GetCartItemInfo(cartID, cartItemID).Return(map[string]string{
							keys.ItemOfferIDMapKey():    offerID,
							keys.ItemCountMapKey():      strconv.FormatUint(uint64(count), 10),
							keys.ItemOfferPriceMapKey(): fmt.Sprintf("%f", offerPrice),
						}).Times(1)
						dbs.EXPECT().GetCartItemOfferTitle(cartID, cartItemID).Return(offerTitle)

						cart, err := s.GetCart(cartID)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(cart.Items)).Should(Equal(1))
						Expect(cart.CartType).Should(Equal(cartType))
						Expect(cart.Items[0].CartItemID).Should(Equal(cartItemID))
						Expect(cart.Items[0].OfferPrice).Should(Equal(offerPrice))
						Expect(cart.Items[0].Count).Should(Equal(count))
						Expect(cart.Items[0].OfferTitle).Should(Equal(offerTitle))
						Expect(cart.Items[0].OfferID).Should(Equal(offerID))
					})
				})
			})
			Describe("Get Cart Type", func() {
				Context("Cart exists", func() {
					It("returns cart type", func() {
						//test data
						const (
							cartID   = "1"
							cartType = api.CartTypeAuthorized
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().GetCartType(cartID).Return(uint8(cartType)).Times(1)

						cartTypeR, err := s.GetCartType(cartID)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(cartTypeR).Should(Equal(cartType))
					})
				})
				Context("Cart does not exist", func() {
					It("replies Not Found", func() {
						//test data
						const (
							cartID = "1"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(false).Times(1)

						_, err := s.GetCartType(cartID)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartNotFound)).Should(BeTrue())
					})
				})
			})
			Describe("Init Cart", func() {
				Context("Cart exists", func() {
					It("replies Already Exists", func() {
						//test data
						const (
							cartID   = "1"
							cartType = api.CartTypeAuthorized
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)

						err := s.InitCart(cartID, cartType)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartAlreadyExists)).Should(BeTrue())
					})
				})
				Context("Cart does not exist", func() {
					It("replies OK", func() {
						//test data
						const (
							cartID   = "1"
							cartType = api.CartTypeAuthorized
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(false).Times(1)
						dbs.EXPECT().AddCartID(cartID).Times(1)
						dbs.EXPECT().SetCartType(cartID, uint8(cartType)).Times(1)
						err := s.InitCart(cartID, cartType)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
			Describe("Cart Exists", func() {
				Context("Cart actually exists", func() {
					It("replies true", func() {
						//test data
						const (
							cartID = "1"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)

						exists, err := s.CartExists(cartID)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).Should(BeTrue())
					})
				})
				Context("Cart does not actually exist", func() {
					It("replies false", func() {
						//test data
						const (
							cartID = "1"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(false).Times(1)

						exists, err := s.CartExists(cartID)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).Should(BeFalse())
					})
				})

			})
			Describe("Merge Carts", func() {
				Context("Both carts exist", func() {
					It("", func() {
						//test data
						const (
							targetCartID     = "1"
							sourceCartID     = "2"
							targetCartType   = api.CartTypeAuthorized
							sourceCartType   = api.CartTypeGuest
							sourceCartItemID = "101"
							sourceOfferPrice = float64(4000)
							sourceCount      = uint(1)
							sourceOfferID    = "10101"
						)
						var sourceOfferTitle = map[string]interface{}{"ru": "Элемент корзины 101"}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(targetCartID).Return(true).Times(2)
						dbs.EXPECT().CartIDIsPresent(sourceCartID).Return(true).Times(2)
						dbs.EXPECT().GetCartType(sourceCartID).Return(uint8(sourceCartType)).Times(1)
						dbs.EXPECT().GetMergedCartIDs(sourceCartID).Return([]string{}).Times(1)
						dbs.EXPECT().GetCartItemIDs(sourceCartID).Return([]string{sourceCartItemID}).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(sourceCartID, sourceCartItemID).Return(true).Times(1)
						dbs.EXPECT().GetCartItemInfo(sourceCartID, sourceCartItemID).Return(map[string]string{
							keys.ItemOfferIDMapKey():    sourceOfferID,
							keys.ItemCountMapKey():      strconv.FormatUint(uint64(sourceCount), 10),
							keys.ItemOfferPriceMapKey(): fmt.Sprintf("%f", sourceOfferPrice),
						}).Times(1)
						dbs.EXPECT().GetCartItemOfferTitle(sourceCartID, sourceCartItemID).Return(sourceOfferTitle)
						dbs.EXPECT().RemoveCartCompletely(sourceCartID).Times(1)
						dbs.EXPECT().AddToMergedCartIDs(targetCartID).Times(1)
						dbs.EXPECT().AddToMergedCartIDs(targetCartID, sourceCartID).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(targetCartID, sourceCartItemID).Return(false).Times(1)
						dbs.EXPECT().AddCartItemID(targetCartID, sourceCartItemID).Times(1)
						dbs.EXPECT().AddCartItemInfo(targetCartID, sourceCartItemID, map[string]interface{}{
							keys.ItemOfferIDMapKey():    sourceOfferID,
							keys.ItemCountMapKey():      sourceCount,
							keys.ItemOfferPriceMapKey(): sourceOfferPrice,
						}).Times(1)
						dbs.EXPECT().AddCartItemOfferTitle(targetCartID, sourceCartItemID, sourceOfferTitle).Times(1)
						err := s.MergeCarts(targetCartID, sourceCartID)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
				Context("Guest cart only", func() {
					It("", func() {
						//test data
						const (
							targetCartID     = "1"
							sourceCartID     = "2"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(targetCartID).Return(false).Times(1)
						err := s.MergeCarts(targetCartID, sourceCartID)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartNotFound)).Should(BeTrue())
					})

				})
				Context("User cart only", func() {
					It("", func() {
						//test data
						const (
							targetCartID     = "1"
							sourceCartID     = "2"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(targetCartID).Return(true).Times(1)
						dbs.EXPECT().CartIDIsPresent(sourceCartID).Return(false).Times(1)
						dbs.EXPECT().AddToMergedCartIDs(targetCartID, sourceCartID).Times(1)
						err := s.MergeCarts(targetCartID, sourceCartID)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
				Context("Both carts do not even exist", func() {
					It("", func() {
						//test data
						const (
							targetCartID     = "1"
							sourceCartID     = "2"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(targetCartID).Return(false).Times(1)
						err := s.MergeCarts(targetCartID, sourceCartID)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartNotFound)).Should(BeTrue())
					})
				})
			})
			Describe("Add Cart Item", func() {
				Context("Invalid cart ID", func() {
					It("", func() {
						//test data
						const (
							cartID   = ""
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}

						err := s.AddCartItem(cartID, cartItem)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.InvalidCartID))
					})
				})
				Context("Invalid cart item ID", func() {
					It("", func() {
						//test data
						const (
							cartID   = "1"
						)
						var cartItem = &api.CartItem{
							CartItemID: "",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						err := s.AddCartItem(cartID, cartItem)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.InvalidCartItemID))
					})
				})
				Context("Cart Item does not exist yet, cart exists", func() {
					It("replies OK", func() {
						//test data
						const (
							cartID   = "1"
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(false).Times(1)
						dbs.EXPECT().AddCartItemID(cartID, cartItem.CartItemID).Times(1)
						dbs.EXPECT().AddCartItemInfo(cartID, cartItem.CartItemID, map[string]interface{}{
							keys.ItemOfferIDMapKey():    cartItem.OfferID,
							keys.ItemCountMapKey():      cartItem.Count,
							keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
						}).Times(1)
						dbs.EXPECT().AddCartItemOfferTitle(cartID, cartItem.CartItemID, cartItem.OfferTitle).Times(1)

						err := s.AddCartItem(cartID, cartItem)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
				Context("Cart Item exists already", func() {
					It("replies CartItem AlreadyExists", func() {
						//test data
						const (
							cartID   = "1"
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(true).Times(1)

						err := s.AddCartItem(cartID, cartItem)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartItemAlreadyExists)).Should(BeTrue())
					})
				})
				Context("Cart does not yet exist", func() {
					It("should not auto init cart, replies Not Found", func() {
						//test data
						const (
							cartID   = "1"
							cartType = api.CartTypeAuthorized
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(false).Times(1)
						//dbs.EXPECT().AddCartID(cartID).Times(1)
						//dbs.EXPECT().SetCartType(cartID, uint8(cartType)).Times(1)
						//dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(false).Times(1)
						//dbs.EXPECT().AddCartItemID(cartID, cartItem.CartItemID).Times(1)
						//dbs.EXPECT().AddCartItemInfo(cartID, cartItem.CartItemID, map[string]interface{}{
						//	keys.ItemOfferIDMapKey():    cartItem.OfferID,
						//	keys.ItemCountMapKey():      cartItem.Count,
						//	keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
						//}).Times(1)
						//dbs.EXPECT().AddCartItemOfferTitle(cartID, cartItem.CartItemID, cartItem.OfferTitle).Times(1)

						err := s.AddCartItem(cartID, cartItem)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartNotFound)).Should(BeTrue())
					})
				})
			})

			Describe("Update Cart Item", func() {
				Context("cart item does not exist", func() {
					It("should reply Not Found", func() {
						//test data
						const (
							cartID   = "1"
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(false).Times(1)

						err := s.UpdateCartItem(cartID, cartItem)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartItemNotFound)).Should(BeTrue())
					})
				})
				Context("cart item exists", func() {
					It("should reply OK", func() {
						//test data
						const (
							cartID   = "1"
						)
						var cartItem = &api.CartItem{
							CartItemID: "101",
							OfferID:    "10101",
							OfferPrice: float64(4000),
							OfferTitle: map[string]interface{}{"ru": "Элемент корзины 101"},
							Count:      uint(1),
						}
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(2)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItem.CartItemID).Return(false).Times(1)
						dbs.EXPECT().RemoveCartItemID(cartID, cartItem.CartItemID).Times(1)
						dbs.EXPECT().RemoveCartItemInfo(cartID, cartItem.CartItemID).Times(1)
						dbs.EXPECT().RemoveCartItemOfferTitle(cartID, cartItem.CartItemID).Times(1)
						dbs.EXPECT().AddCartItemID(cartID, cartItem.CartItemID).Times(1)
						dbs.EXPECT().AddCartItemInfo(cartID, cartItem.CartItemID, map[string]interface{}{
							keys.ItemOfferIDMapKey():    cartItem.OfferID,
							keys.ItemCountMapKey():      cartItem.Count,
							keys.ItemOfferPriceMapKey(): cartItem.OfferPrice,
						}).Times(1)
						dbs.EXPECT().AddCartItemOfferTitle(cartID, cartItem.CartItemID, cartItem.OfferTitle).Times(1)

						err := s.UpdateCartItem(cartID, cartItem)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
			Describe("Remove Cart Item", func() {
				Context("cart item does not exist", func() {
					It("should reply Not Found", func() {
						//test data
						const (
							cartID     = "1"
							cartItemID = "101"
						)
						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItemID).Return(false).Times(1)

						err := s.RemoveCartItem(cartID, cartItemID)
						Expect(err).Should(HaveOccurred())
						Expect(err.(carterrors.CartError).IsType(carterrors.CartItemNotFound)).Should(BeTrue())
					})
				})
				Context("cart item exists", func() {
					It("should reply OK", func() {
						//test data
						const (
							cartID     = "1"
							cartItemID = "101"
						)

						//expectations on DB server mock
						dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
						dbs.EXPECT().CartItemIDIsPresent(cartID, cartItemID).Return(true).Times(1)
						dbs.EXPECT().RemoveCartItemID(cartID, cartItemID).Times(1)
						dbs.EXPECT().RemoveCartItemInfo(cartID, cartItemID).Times(1)
						dbs.EXPECT().RemoveCartItemOfferTitle(cartID, cartItemID).Times(1)

						err := s.RemoveCartItem(cartID, cartItemID)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})
	})
})
