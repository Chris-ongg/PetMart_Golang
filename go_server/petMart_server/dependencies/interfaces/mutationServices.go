package interfaces

import (
	"go_server/petMart_server/model"
)

type MutationServices interface {
	RegisterNewCustomer(name string, email string , password string) error
	LogoutCustomer(email string) string
	RegisterNewPets(newPet *model.PetDetails) error
	SaveShoppingCart(currentCart *model.CustomerCart) error
	SaveOrder(purchasedCart *model.CustomerCart) error
}
