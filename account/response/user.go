package response

import "account/models"

type AddUser struct {
	User             *models.User
	ErrValidation    error
	Err              string
	ErrorUniqueField []ErrorUniqueField
}

type LoginResponse struct {
	Token         string
	Err           string
	ErrValidation error
}

type ValidateToken struct {
	User          *models.User
	Err           string
	ErrValidation error
}

type UpdatePassword struct {
	Err       string
	IsUpdated bool

	ErrValidation error
}

type UpdateUser struct {
	Err              string
	ErrValidation    error
	User             *models.User
	ErrorUniqueField []ErrorUniqueField
}

type GetUser struct {
	User models.User
	Err string
}

type AddRoleToUser struct {
	User *models.User
	Err string
	ErrorUniqueField []ErrorUniqueField
}

type DeleteUser struct {
	Err       string
	IsDeleted bool
}
