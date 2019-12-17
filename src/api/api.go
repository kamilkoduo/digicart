package api

type CartType uint8

const (
	CartType_Authorized CartType = 0
	CartType_Guest      CartType = 1
)

type Cart struct {
	CartID        string      `jsonapi:"primary, cart_id"`
	CartType      CartType    `jsonapi:"attr, cart_type"`
	MergedCartIDs []string    `jsonapi:"attr, merged_cart_ids"`
	Items         []*CartItem `jsonapi:"relation, cart_items"`
}

type CartItem struct {
	CartItemID string            `jsonapi:"primary,cart_item"`
	OfferID    string            `jsonapi:"attr,offer_id"`
	OfferPrice float64           `jsonapi:"attr,offer_price"`
	OfferTitle map[string]string `jsonapi:"attr,offer_title"`
	Count      int               `jsonapi:"attr,count"`
}

type CartApi interface {
	GetCart(cartID string, cartType CartType) (*Cart, error)
	MergeCarts(TargetCartID string, SourceCartID string) error
	AddCartItem(cartID string) error
	UpdateCartItem(cartID string) error
	RemoveCartItem(cartID string) error
}
