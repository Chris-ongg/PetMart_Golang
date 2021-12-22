package database

import "go_server/petMart_server/model"

type CustomerAccount struct {
	Name     string
	EmailAdd string
	Password string
	Registered bool
	AccessToken string
	TokenExpiry int64
}

type CustomerPets struct {
	Name string
	Gender string
	PetBreed string
	Species string
	Weight string
	OwnerEmail string
	HealthConcern string
	Age string
}

type CustomerActivity struct{
	Email string
	EmptyCart bool
	CartItems []*model.CartItemsInput
	PastOrders []*model.PastOrders
}

type OrderID struct {
	Name     string
	LatestID int
}

