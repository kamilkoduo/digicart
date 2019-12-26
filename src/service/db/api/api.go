package api

type CartDBAPI interface {
	CartIDIsPresent(cartID string) bool
	CartItemIDIsPresent(cartID, cartItemID string) bool
	AddCartID(cartID string)
	SetCartType(cartID string, cartType uint8)
	GetMergedCartIDs(cartID string) []string
	AddCartItemID(cartID, cartItemID string)
	RemoveCartItemID(cartID, cartItemID string)
	AddCartItemInfo(cartID, cartItemID string, cartItemInfo map[string]interface{})
	RemoveCartItemInfo(cartID, cartItemID string)
	AddCartItemOfferTitle(cartID, cartItemID string, cartItemOfferTitle map[string]interface{})
	RemoveCartItemOfferTitle(cartID, cartItemID string)
	AddToMergedCartIDs(cartID string, mergedID ...string)
	GetCartType(cartID string) uint8
	GetCartItemInfo(cartID, cartItemID string) map[string]string
	GetCartItemOfferTitle(cartID, cartItemID string) map[string]interface{}
	GetCartItemIDs(cartID string) []string
	RemoveCartCompletely(cartID string)
}
