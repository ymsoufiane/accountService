package services

import (
	"account/context"
	"account/models"
	"account/request"
	"account/response"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var notExist uint

type UserService interface {
	GenerateToken(user models.User) (string, error)
	ValidateUser(user *models.User) (bool, error)
	ValidateToken(tokenString string) (*models.User, error)
	UpdatePassword(updatePassword request.UpdatePassword, user models.User) error
	checkUniqueField(user models.User, ignoreId uint) (bool, []response.ErrorUniqueField)
	AddUser(user *models.User) (int, response.AddUser)
	UpdateUser(user *models.User) (int, response.UpdateUser)
	AddRoleToUser(user *models.User, role models.Role) (int, response.AddRoleToUser)
	AddRolesToUser(user *models.User, roles []models.Role) (int, response.AddRoleToUser)
}
type UserServiceImpl struct {
}

var userRep *models.UserRepository
var roleService RoleService

func init() {
	notExist = 0
	userRep = models.UserRep
	roleService = RoleServiceImpl{}
}

func (u UserServiceImpl) AddUser(user *models.User) (int, response.AddUser) {
	ok, err1 := u.checkUniqueField(*user, notExist)
	if !ok {
		return http.StatusBadRequest, response.AddUser{ErrorUniqueField: err1}
	}

	roleValidation := roleService.ValidateRoles(user.Roles)
	if !roleValidation.IsValid {
		return http.StatusBadRequest, response.AddUser{ErrorUniqueField: roleValidation.ErrorValidation}
	}

	userRes, err2 := userRep.Create(user)
	if err2 != nil {
		return http.StatusInternalServerError, response.AddUser{Err: response.MessageErrInternal}
	}

	return http.StatusOK, response.AddUser{User: userRes}
}

func (u UserServiceImpl) AddRoleToUser(user *models.User, role models.Role) (int, response.AddRoleToUser) {
	result := make(chan ValidateRole)
	go roleService.ValidateRole(role, result)
	tempResult := <-result
	if tempResult.IsValid {
		user.Roles = append(user.Roles, role)
		userRep.Save(user)
		return http.StatusOK, response.AddRoleToUser{User: user}
	}
	return http.StatusBadRequest, response.AddRoleToUser{ErrorUniqueField: tempResult.ErrorValidation}

}

func (u UserServiceImpl) AddRolesToUser(user *models.User, roles []models.Role) (int, response.AddRoleToUser) {
	result := roleService.ValidateRoles(roles)
	if result.IsValid {
		user.Roles = append(user.Roles, roles...)
		userRep.Save(user)
		return http.StatusOK, response.AddRoleToUser{User: user}
	}
	return http.StatusBadRequest, response.AddRoleToUser{ErrorUniqueField: result.ErrorValidation}
}

func (u UserServiceImpl) UpdateUser(user *models.User) (int, response.UpdateUser) {
	ok, err1 := u.checkUniqueField(*user, user.ID)
	if !ok {
		return http.StatusBadRequest, response.UpdateUser{ErrorUniqueField: err1}
	}

	roleValidation := roleService.ValidateRoles(user.Roles)
	if !roleValidation.IsValid {
		return http.StatusBadRequest, response.UpdateUser{ErrorUniqueField: roleValidation.ErrorValidation}
	}

	userRes, err2 := userRep.Update(user)
	if err2 != nil {
		return http.StatusInternalServerError, response.UpdateUser{Err: response.MessageErrInternal}
	}
	return http.StatusOK, response.UpdateUser{User: userRes}
}

func (u UserServiceImpl) GenerateToken(user models.User) (string, error) {

	keyString := context.Config.Jwt.Secret
	keyBinary := []byte(keyString)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["ExpiresAt"] = time.Now().Add(45 * time.Minute)
	claims["username"] = user.Username
	claims["firstName"] = user.FirstName
	claims["id"] = user.ID
	claims["lastName"] = user.LastName
	claims["email"] = user.Email
	claims["phoneNumber"] = user.PhoneNumber
	claims["prefixePhoneNumber"] = user.PrefixePhoneNumber
	claims["roles"] = user.Roles

	return token.SignedString(keyBinary)

}

func (u UserServiceImpl) ValidateUser(user *models.User) (bool, error) {
	username := user.Username
	userResult := models.UserRep.FindByUserName(user.Username)

	if userResult.ID == notExist {
		context.WarningLogger.Println("user " + username + " not exist ")
		return false, errors.New("user " + username + " not exist ")
	}

	err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(user.Password))

	isMatch := err == nil
	if isMatch {
		*user = *userResult
		return true, nil
	} else {
		context.WarningLogger.Println(err)
	}
	return false, errors.New("password incorrect !")

}

func (u UserServiceImpl) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		keyString := context.Config.Jwt.Secret
		keyBinary := []byte(keyString)
		return keyBinary, nil
	})
	if err != nil {
		context.WarningLogger.Println(err)
		return nil, err
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("error convet claim to map ")
	}

	var user models.User
	stringify, _ := json.Marshal(claim)
	json.Unmarshal(stringify, &user)

	return &user, nil

}

func (u UserServiceImpl) UpdatePassword(updatePassword request.UpdatePassword, user models.User) error {
	if updatePassword.NewPassword != updatePassword.ConfirmPassword {
		return errors.New("New Password field did not match your input under Confirm New Password,")
	}
	user = *models.UserRep.FindById(int(user.ID))
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updatePassword.OldPassword))
	NotMatch := err != nil
	if NotMatch {
		return errors.New("old password incorrect !!")
	}
	user.Password = updatePassword.NewPassword
	models.UserRep.UpdatePassword(&user)
	return nil

}

func (u UserServiceImpl) checkUniqueField(user models.User, ignoreId uint) (bool, []response.ErrorUniqueField) {
	result := models.UserRep.FindByEmail(user.Email)
	var errors []response.ErrorUniqueField
	if result.ID != ignoreId {
		errors = append(errors, response.ErrorUniqueField{Name: "email", Message: "this email already in use"})
	}
	result = models.UserRep.FindByPhoneNumber(user.PhoneNumber)
	if result.ID != ignoreId {
		errors = append(errors, response.ErrorUniqueField{Name: "phoneNumber", Message: "this phone number already in use"})
	}
	result = models.UserRep.FindByUserName(user.Username)
	if result.ID != ignoreId {
		errors = append(errors, response.ErrorUniqueField{Name: "userName", Message: "this username already in use"})
	}
	if len(errors) > 0 {
		return false, errors
	}
	return true, nil

}
