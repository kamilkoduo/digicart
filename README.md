# Diginavis Cart 

- Test project for Diginavis

##### API
- GET /api/v1/cart-api/my  
response: cart
- POST /api/v1/cart-api/my/items/:id

- DELETE /api/v1/cart-api/my/items/:id
- PUT /api/v1/cart-api/my/items/:id
  
##### Database:
- Redis
  
##### Cart item fields:
- offer_id - string
- offer_price - number
- offer_title - map[string]string{"en": "title" }
- count - int
  
##### Request/Response format: 
- JSONAPI `google/jsonapi`
  
##### Authentication headers:
- X-User-ID
- X-Guest-ID