package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go_server/petMart_server/model"
)

func (r *mutationResolver) SaveShoppingCart(ctx context.Context, saveCart *model.CustomerCart) (*model.CustomerEmailAdd, error) {
	status := r.MutationService.SaveShoppingCart(saveCart)
	if status != nil {
		fmt.Println(status)
	}
	customerEmail := model.CustomerEmailAdd{Email: saveCart.Email}
	return &customerEmail, status
}

func (r *mutationResolver) SaveCustomerOrder(ctx context.Context, saveOrder *model.CustomerCart) (*model.CustomerEmailAdd, error) {
	status := r.MutationService.SaveOrder(saveOrder)
	if status != nil {
		fmt.Println(status)
	}
	customerEmail := model.CustomerEmailAdd{Email: saveOrder.Email}
	return &customerEmail, status
}

func (r *queryResolver) SearchShoppingCart(ctx context.Context, cart *model.CustomerEmailInput) ([]*model.CartItems, error) {
	cartResult, err := r.QueryService.GetShoppingCart(cart.Email)
	var tempArr []*model.CartItems
	for _, item := range cartResult {
		tempItem := model.CartItems{
			ItemID:    item.ItemID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			ImagePath: item.ImagePath,
		}
		tempArr = append(tempArr, &tempItem)
	}
	return tempArr, err
}

func (r *queryResolver) CustomerPastOrders(ctx context.Context, customer *model.CustomerEmailInput) ([]*model.PastOrders, error) {
	result, err := r.QueryService.GetCustomerPastOrders(customer.Email)
	return result, err
}
