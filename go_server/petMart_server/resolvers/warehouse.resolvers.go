package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go_server/petMart_server/model"
)

func (r *queryResolver) SearchWarehouse(ctx context.Context, search model.SearchInput) ([]*model.ProductList, error) {
	retrievedProducts, err := r.QueryService.GetProductsFromWarehouse(search.SearchType, search.HealthConcern, search.Species)
	return retrievedProducts, err
}
