package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"go_server/petMart_server/dependencies/interfaces"
)

type Resolver struct{
	QueryService interfaces.QueryServices
	MutationService interfaces.MutationServices
}
