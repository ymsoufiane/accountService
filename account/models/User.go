package models

import (
	"account/context"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string `gorm:"unique;not null;type:varchar(100);"  `
	LastName           string `gorm:"not null;type:varchar(100);"`
	FirstName          string `gorm:"not null;type:varchar(100);" `
	Email              string `gorm:"unique;not null;type:varchar(100);" `
	Password           string `gorm:"not null;type:varchar(255);" `
	PhoneNumber        string `gorm:"unique;type:varchar(20);"`
	PrefixePhoneNumber string `gorm:"type:varchar(5);" `
	Roles              []Role `gorm:"many2many:user_roles;"`
}

func (user *User) FromRequest(userRequest interface{}) *User {
	var userModel User
	ByteUser, err := json.Marshal(userRequest)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(ByteUser, &userModel)
	if err != nil {
		return nil
	}
	return &userModel

}

type UserRepository struct {
	db *gorm.DB
}

var UserRep *UserRepository

func init() {
	db := context.Db
	UserRep = &UserRepository{db}
	err := db.AutoMigrate(&User{})
	if err != nil {
		context.WarningLogger.Println(err)
		fmt.Println(err)
		fmt.Println("error Migrate User table ")
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	err := r.db.Create(user).Error

	return user, err
}

func (r *UserRepository) Save(user *User) (*User, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	err := r.db.Save(user).Error
	return user, err
}

func (r *UserRepository) Update(user *User) (*User, error) {
	user.Password = ""
	err := r.db.Updates(user).Error
	return user, err
}

func (r *UserRepository) UpdatePassword(user *User) (*User, error) {

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	err := r.db.Model(&user).Update("password", user.Password).Error
	return user, err
}

func (r *UserRepository) Delete(user *User) (*User, error) {
	err := r.db.Delete(user).Error
	return user, err
}

func (r *UserRepository) DeleteById(id int) error {
	err := r.db.Delete(&User{}, id).Error
	return err
}

func (r *UserRepository) FindById(id int) *User {
	var user User

	r.db.First(&user, id)

	return &user
}

func (r *UserRepository) FindByUserName(username string) *User {
	var user User
	r.db.Where("username=?", username).Preload("Roles").Preload("Roles.Priviliges").First(&user)
	return &user
}

func (r *UserRepository) FindByPhoneNumber(phoneNumber string) *User {
	var user User
	r.db.Where("phone_number=?", phoneNumber).Preload("Roles").Preload("Roles.Priviliges").First(&user)
	return &user
}

func (r *UserRepository) FindByEmail(email string) *User {
	var user User
	r.db.Where("email=?", email).Preload("Roles").Preload("Roles.Priviliges").First(&user)
	return &user
}

func (r *UserRepository) GetAll()[]User{
	var users []User
	r.db.Find(&users)
	return users
}