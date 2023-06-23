package request

import "account/models"

type AddRole struct {
	RoleName    string             `json:"roleName" valid:"required"`
	Description string             `json:"description" valid:"required"`
	Priviliges  []models.Privilige `valid:"-" `
}
