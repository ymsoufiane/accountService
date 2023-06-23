package services

import (
	"account/models"
	"account/response"
	"errors"
)


var priviligeRep models.PriviligeRepository

func init(){
	priviligeRep=*models.PriviligeRep
}

type PriviligeService interface {
	ValidatePrivilige(priv *models.Privilige)(bool,error)
	ValidatePriviliges(privs []models.Privilige)(bool,[]response.ErrorUniqueField)
}

type PriviligeServiceImpl struct{}

func (s PriviligeServiceImpl) ValidatePrivilige(priv *models.Privilige)(bool,error){
	id := priv.ID
	result :=priviligeRep.FindById(int(id))
	if result.ID!=notExist{
		return true,nil
	}
	result = priviligeRep.FindPriviligeByName(&priv.Privilige)
	if result.ID!=notExist{
		return false,errors.New("this privilige name already in use")
	}
	return true,nil
}
func (s PriviligeServiceImpl) ValidatePriviliges(privs []models.Privilige)(bool,[]response.ErrorUniqueField){
	isValid :=	true
	var ok bool
	var err error
	var result []response.ErrorUniqueField
	for _,priv := range privs{
		ok,err=s.ValidatePrivilige(&priv)
		if !ok{
			isValid=false
			result = append(result,response.ErrorUniqueField{Name: priv.Privilige,Message: err.Error()} )
		}
	}
	return isValid,result
}