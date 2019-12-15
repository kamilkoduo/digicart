package api

type Cart struct {
	CartID  string      `jsonapi:"primary, cart"`
	UserID  string      `jsonapi:"attr, user_id"`
	GuestID string      `jsonapi:"attr, guest_id"`
	Items   []*CartItem `jsonapi:"relation, cart_item"`
}

type CartItem struct {
	CartItemID string            `jsonapi:"primary,cart_item"`
	OfferID    string            `jsonapi:"attr,offer_id"`
	OfferPrice float64           `jsonapi:"attr,offer_price"`
	OfferTitle map[string]string `jsonapi:"attr,offer_title"`
	Count      int               `jsonapi:"attr,count"`
}

type CartApi interface {
	GetCart(UserID string, GuestID string) (*Cart,error)
	AddCartItem(UserID string, GuestID string) (error)
	UpdateCartItem(UserID string, GuestID string) (error)
	RemoveCartItem(UserID string, GuestID string) (error)
}