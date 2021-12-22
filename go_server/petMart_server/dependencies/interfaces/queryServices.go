package interfaces

import (
	"go_server/petMart_server/model"
)

type QueryServices interface {
	 CustomerPets(email string) ([]*model.PetDetails , error)
	 EmailLogin(email string ,password string) (string , string)
	 GoogleLogin(accessToken string) (string , string , string)
	 GetShoppingCart(email string) ([]*model.CartItemsInput , error)
	 TokenVerification(token string) (*model.AuthenticatedUser , error)
	 GetCustomerPastOrders(email string) ([]*model.PastOrders ,  error)
	 GetProductsFromWarehouse(searchType int , health string, species string) ([]*model.ProductList , error)
}

