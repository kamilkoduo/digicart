//+build unit

package service_test

//
//func setupAPIServerWithMock(t gomock.TestReporter) *service.CartAPIServer {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	MockCartDBAPI := mocks.NewMockCartDBAPI(mockCtrl)
//
//	// lets say we mock a redis DB with:
//	// cart: id "1", type user ("0")
//	// cart item: id "101"
//
//	MockCartDBAPI.EXPECT().CartIDIsPresent(cartID).Return(true)
//	//MockCartDBAPI.EXPECT().GetCartItemIDs(cartID).Return([]string{cartItemID})
//	//MockCartDBAPI.EXPECT().CartItemIDIsPresent(cartID, cartItemID).Return(true)
//	//MockCartDBAPI.EXPECT().GetCartType(cartID).Return(uint8(0))
//	/*	MockCartDBAPI.EXPECT().AddCartID()
//		MockCartDBAPI.EXPECT().SetCartType(cartID string, cartType uint8)
//		MockCartDBAPI.EXPECT().GetMergedCartIDs(cartID string) []string
//		MockCartDBAPI.EXPECT().AddCartItemID(cartID, cartItemID string)
//		MockCartDBAPI.EXPECT().RemoveCartItemID(cartID, cartItemID string)
//		MockCartDBAPI.EXPECT().AddCartItemInfo(cartID, cartItemID string, cartItemInfo map[string]interface{})
//		MockCartDBAPI.EXPECT().RemoveCartItemInfo(cartID, cartItemID string)
//		MockCartDBAPI.EXPECT().AddCartItemOfferTitle(cartID, cartItemID string, cartItemOfferTitle map[string]interface{})
//		MockCartDBAPI.EXPECT().RemoveCartItemOfferTitle(cartID, cartItemID string)
//		MockCartDBAPI.EXPECT().AddToMergedCartIDs(cartID string, mergedID ...string)
//		MockCartDBAPI.EXPECT().GetCartType(cartID string) uint8
//		MockCartDBAPI.EXPECT().GetCartItemInfo(cartID, cartItemID string) map[string]string
//		MockCartDBAPI.EXPECT().GetCartItemOfferTitle(cartID, cartItemID string) map[string]interface{}
//		MockCartDBAPI.EXPECT().GetCartItemIDs(cartID string) []string
//		MockCartDBAPI.EXPECT().RemoveCartCompletely(cartID string)
//	*/
//
//	s := &service.CartAPIServer{}
//	s.RegisterCartDBAPIServer(MockCartDBAPI)
//	return s
//}
