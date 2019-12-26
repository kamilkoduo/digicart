package keys

import (
	"strings"
)

const (
	sep            = ":"
	cartSet        = "carts"
	cartType       = "type"
	mergedSet      = "merged"
	items          = "items"
	itemInfo       = "info"
	itemOfferTitle = "offer_title"
	itemOfferID    = "offer_id"
	itemOfferPrice = "offer_price"
	itemCount      = "count"
)

func CartSet() string {
	return cartSet
}
func CartPrefix(cartID string) string {
	return strings.Join([]string{cartSet, cartID}, sep)
}
func CartTypeKey(cartID string) string {
	return strings.Join([]string{cartSet, cartID, cartType}, sep)
}
func CartMergedSet(cartID string) string {
	return strings.Join([]string{cartSet, cartID, mergedSet}, sep)
}
func CartItemSet(cartID string) string {
	return strings.Join([]string{cartSet, cartID, items}, sep)
}
func ItemInfoMap(cartID, cartItemID string) string {
	return strings.Join([]string{cartSet, cartID, items, cartItemID, itemInfo}, sep)
}
func ItemOfferTitleMap(cartID, cartItemID string) string {
	return strings.Join([]string{cartSet, cartID, items, cartItemID, itemOfferTitle}, sep)
}

func ItemOfferIDMapKey() string {
	return itemOfferID
}
func ItemOfferPriceMapKey() string {
	return itemOfferPrice
}
func ItemCountMapKey() string {
	return itemCount
}
