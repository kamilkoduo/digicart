//+build unit

package service_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/kamilkoduo/digicart/src/api"
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
	BeforeSuite(func() {
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
		Context("cart exists with one item", func() {
			Context("check cart exists", func() {
				It("should be correct", func() {
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
					dbs.EXPECT().CartIDIsPresent(cartID).Return(true).Times(1)
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
					Expect(cart.CartType).Should(Equal(api.CartTypeAuthorized))
					Expect(cart.Items[0].CartItemID).Should(Equal(cartItemID))
					Expect(cart.Items[0].OfferPrice).Should(Equal(offerPrice))
					Expect(cart.Items[0].Count).Should(Equal(count))
					Expect(cart.Items[0].OfferTitle).Should(Equal(offerTitle))
					Expect(cart.Items[0].OfferID).Should(Equal(offerID))

				})
			})
		})
	})
})
