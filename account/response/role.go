package response

import "account/models"

type AddRole struct {
	Role             *models.Role
	ErrValidation    error
	Err              string
	ErrorUniqueField []ErrorUniqueField
}

type GetRole struct {
	Role *models.Role
	Err  string
}
