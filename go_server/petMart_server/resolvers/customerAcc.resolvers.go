package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go_server/petMart_server/middleware"
	"go_server/petMart_server/model"
	"os"
)

func (r *mutationResolver) CustomerRegistration(ctx context.Context, register model.CustomerRegistrationInput) (*model.CustomerDetails, error) {
	newCustomer := model.CustomerDetails{
		Name:  register.Name,
		Email: register.Email,
	}

	status := r.MutationService.RegisterNewCustomer(register.Name, register.Email, register.Password)

	if status != nil {
		fmt.Println(status)
	}

	return &newCustomer, nil
}

func (r *mutationResolver) CustomerLogout(ctx context.Context, logout model.CustomerEmailInput) (*model.CustomerDetails, error) {
	customerName := r.MutationService.LogoutCustomer(logout.Email)

	writeCookie := middleware.GetCookieAccess(ctx)
	writeCookie.SetToken("")

	customerDetails := model.CustomerDetails{
		Name:  customerName,
		Email: logout.Email,
	}

	return &customerDetails, nil
}

func (r *queryResolver) CustomerEmailLogin(ctx context.Context, login *model.CustomerLoginInput) (*model.AuthenticatedUser, error) {
	status, name := r.QueryService.EmailLogin(login.Email, login.Password)
	var authenticatedUser model.AuthenticatedUser
	if status == "fail" {
		authenticatedUser = model.AuthenticatedUser{
			Name:        "",
			Email:       "",
			AccessToken: "",
		}
	} else {
		authenticatedUser = model.AuthenticatedUser{
			Name:        name,
			Email:       login.Email,
			AccessToken: status,
		}
		writeCookie := middleware.GetCookieAccess(ctx)
		writeCookie.SetToken(status)
	}
	return &authenticatedUser, nil
}

func (r *queryResolver) CustomerGoogleLogin(ctx context.Context, login model.AccessTokenInput) (*model.AuthenticatedUser, error) {
	name , email , token := r.QueryService.GoogleLogin(login.Token)
	var authenticatedUser model.AuthenticatedUser
	authenticatedUser = model.AuthenticatedUser{
			Name:        name,
			Email:       email,
			AccessToken: token,
		}
		writeCookie := middleware.GetCookieAccess(ctx)
		writeCookie.SetToken(token)
	return &authenticatedUser, nil
}

func (r *queryResolver) VerifyAccessToken(ctx context.Context, token model.AccessTokenInput) (*model.AuthenticatedUser, error) {
	user, err := r.QueryService.TokenVerification(token.Token)
	return user, err
}

func (r *queryResolver) EncryptionKey(ctx context.Context) (*model.EncryptKey, error) {
	key := model.EncryptKey{
		Key: os.Getenv("SECRETKEY"),
	}
	return &key, nil
}
