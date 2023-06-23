package models

import (
	"account/context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleName    string      `gorm:"unique;not null;type:varchar(100);"`
	Description string      `gorm:"default:null;type:text;"`
	Users       []User      `gorm:"many2many:user_roles;"`
	Priviliges  []Privilige `gorm:"many2many:role_priviliges;"`
}

func (role *Role) FromRequest(roleRequest interface{}) *Role {
	var roleModel Role
	ByteRole, err := json.Marshal(roleRequest)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(ByteRole, &roleModel)
	if err != nil {
		return nil
	}
	return &roleModel

}

type RoleRepository struct {
	db *gorm.DB
}

var RoleRep *RoleRepository
var notExist uint

func init() {
	notExist = 0
	db := context.Db
	RoleRep = &RoleRepository{db}
	err := db.AutoMigrate(&Role{})
	if err != nil {
		context.WarningLogger.Println(err)
		fmt.Println(err)
		fmt.Println("error Migrate Role table  ")
	}
}

func (r *RoleRepository) Create(role *Role) (*Role, error) {
	err := RoleRep.db.Create(role).Error

	return role, err
}

func (r *RoleRepository) Save(role *Role) (*Role, error) {
	err := RoleRep.db.Save(role).Error
	return role, err
}

func (r *RoleRepository) Update(role *Role) (*Role, error) {
	err := RoleRep.db.Updates(role).Error
	return role, err
}

func (r *RoleRepository) Delete(role *Role) (*Role, error) {
	err := RoleRep.db.Delete(role).Error
	return role, err
}

func (r *RoleRepository) FindByRoleName(rolename string) (*Role, error) {
	var role Role
	err := r.db.Where("role_name=?", rolename).First(&role).Error
	return &role, err

}

func (r *RoleRepository) FindById(id int) *Role {
	var role Role
	r.db.Find(&role, id)

	return &role
}
