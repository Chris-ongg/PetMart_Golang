package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go_server/petMart_server/generated"
	"go_server/petMart_server/middleware"
	"go_server/petMart_server/model"
)

func (r *mutationResolver) RegisterPet(ctx context.Context, newPet *model.NewPetDetails) (*model.PetDetails, error) {
	newCustomerPet := model.PetDetails{
		Name:          newPet.Name,
		Gender:        newPet.Gender,
		PetBreed:      newPet.PetBreed,
		Species:       newPet.Species,
		Weight:        newPet.Weight,
		OwnerEmail:    newPet.OwnerEmail,
		HealthConcern: newPet.HealthConcern,
		Age:           newPet.Weight,
	}

	status := r.MutationService.RegisterNewPets(&newCustomerPet)
	return &newCustomerPet, status
}

func (r *queryResolver) GetCustomerPets(ctx context.Context, customer model.CustomerEmailInput) ([]*model.PetDetails, error) {
	var pets []*model.PetDetails
	writeCookie := middleware.GetCookieAccess(ctx)
	writeCookie.SetToken("testing")
	pets, _ = r.QueryService.CustomerPets(customer.Email)
	if pets == nil {
		pets = []*model.PetDetails{}
	}
	//fmt.Println(len(pets))
	//fmt.Printf("%+v" , pets)
	return pets, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
