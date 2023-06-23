package services

import (
	"account/context"
	"account/models"
	"account/response"
	"net/http"
)

var roleRep models.RoleRepository
var priviligeService PriviligeService

func init() {
	roleRep = *models.RoleRep
	priviligeService = PriviligeServiceImpl{}
}

type RoleService interface {
	ValidateRoles(roles []models.Role) ValidateRole
	ValidateRole(role models.Role, rsp chan ValidateRole)
	AddRole(role *models.Role) (int, response.AddRole)
	UpdateRole(role *models.Role) (int, response.AddRole)
}

type RoleServiceImpl struct{}

type ValidateRole struct {
	IsValid         bool
	ErrorValidation []response.ErrorUniqueField
}

func (s RoleServiceImpl) AddRole(role *models.Role) (int, response.AddRole) {
	result := make(chan ValidateRole)
	go roleService.ValidateRole(*role, result)
	tempResult := <-result
	if !tempResult.IsValid {
		return http.StatusBadRequest, response.AddRole{ErrorUniqueField: tempResult.ErrorValidation}
	}
	roleRes, err := roleRep.Create(role)
	if err != nil {
		context.WarningLogger.Println(err)
		return http.StatusInternalServerError, response.AddRole{Err: response.MessageErrInternal}
	}

	return http.StatusOK, response.AddRole{Role: roleRes}
}

func (s RoleServiceImpl) UpdateRole(role *models.Role) (int, response.AddRole) {
	result := make(chan ValidateRole)
	go roleService.ValidateRole(*role, result)
	tempResult := <-result
	if !tempResult.IsValid {
		return http.StatusBadRequest, response.AddRole{ErrorUniqueField: tempResult.ErrorValidation}
	}
	roleRes, err := roleRep.Update(role)
	if err != nil {
		context.WarningLogger.Println(err)
		return http.StatusInternalServerError, response.AddRole{Err: response.MessageErrInternal}
	}

	return http.StatusOK, response.AddRole{Role: roleRes}
}

func (s RoleServiceImpl) ValidateRole(role models.Role, rsp chan ValidateRole) {

	id := role.ID
	roleRes := roleRep.FindById(int(id))
	if roleRes.ID != notExist {
		//validate priviliges if role already exist
		isValid, err := priviligeService.ValidatePriviliges(role.Priviliges)
		rsp <- ValidateRole{IsValid: isValid, ErrorValidation: err}
	}

	result, _ := roleRep.FindByRoleName(role.RoleName)
	if result.ID != notExist {
		var err []response.ErrorUniqueField
		err = append(err, response.ErrorUniqueField{Name: role.RoleName, Message: role.RoleName + " already in use"})
		rsp <- ValidateRole{IsValid: false, ErrorValidation: err}
	}
	isValid, err := priviligeService.ValidatePriviliges(role.Priviliges)
	rsp <- ValidateRole{IsValid: isValid, ErrorValidation: err}

}

func (s RoleServiceImpl) ValidateRoles(roles []models.Role) ValidateRole {
	validation := ValidateRole{IsValid: true}
	size := len(roles)
	result := make(chan ValidateRole)

	for i := 0; i < size; i++ {
		go s.ValidateRole(roles[i], result)
	}
	var tempResult ValidateRole
	for i := 0; i < size; i++ {
		tempResult = <-result
		validation.IsValid = tempResult.IsValid && validation.IsValid
		validation.ErrorValidation = append(validation.ErrorValidation, tempResult.ErrorValidation...)
	}
	return validation
}
