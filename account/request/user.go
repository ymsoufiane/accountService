package request

import "account/models"

type TokenValidation struct {
	Token string `json:"token" valid:"required" `
}

type Login struct {
	Username string `valid:"required" json:"username"`
	Password string `valid:"required" json:"password"`
}

type Registre struct {
	FirstName          string        `valid:"required" json:"firstName"`
	LastName           string        `valid:"required" json:"lastName" `
	PhoneNumber        string        `valid:"required" json:"phoneNumber" `
	PrefixePhoneNumber string        `json:"prefixePhoneNumber" valid:"-" `
	Username           string        `valid:"required" json:"username" `
	Email              string        `valid:"required,email" json:"email"`
	Password           string        `valid:"required" json:"password" `
	Roles              []models.Role `valid:"-" `
}

type Update struct {
	FirstName          string `valid:"required" json:"firstName"`
	LastName           string `valid:"required" json:"lastName" `
	PhoneNumber        string `valid:"required" json:"phoneNumber" `
	PrefixePhoneNumber string `json:"prefixePhoneNumber" `
	Username           string `valid:"required" json:"username" `
	Email              string `valid:"required" json:"password" `
	Roles              []models.Role `valid:"-" `
}
type UpdatePassword struct {
	OldPassword     string `valid:"required" json:"oldPassword" `
	NewPassword     string `valid:"required" json:"newPassword" `
	ConfirmPassword string `valid:"required" json:"confirmPassword" `
}
