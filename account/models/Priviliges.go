package models

import (
	//"gorm.io/driver/postgres"

	"account/context"
	"fmt"

	"gorm.io/gorm"
)

type Privilige struct {
	gorm.Model
	Privilige   string `gorm:"unique;not null;type:varchar(255);"  `
	Description string `gorm:"default:null;type:text;"`
	Roles       []Role `gorm:"many2many:role_priviliges;"`
}

type PriviligeRepository struct {
	db *gorm.DB
}

var PriviligeRep *PriviligeRepository

func init() {
	db := context.Db
	PriviligeRep = &PriviligeRepository{db}
	err := db.AutoMigrate(&Privilige{})
	if err != nil {
		context.WarningLogger.Println(err)
		fmt.Println(err)
		fmt.Println("error Migrate Privilige table ")

	}
}

func (p *PriviligeRepository) Create(privilige *Privilige) (*Privilige, error) {
	err := p.db.Create(privilige).Error
	return privilige, err
}

func (p *PriviligeRepository) Save(privilige *Privilige) (*Privilige, error) {
	err := p.db.Save(privilige).Error
	return privilige, err
}

func (p *PriviligeRepository) Update(privilige *Privilige) (*Privilige, error) {
	err := p.db.Updates(privilige).Error
	return privilige, err
}

func (p *PriviligeRepository) Delete(privilige *Privilige) (*Privilige, error) {
	err := p.db.Delete(privilige).Error
	return privilige, err
}

func (p *PriviligeRepository) FindById(id int) *Privilige {
	var privilige Privilige
	p.db.Find(&privilige, id)

	return &privilige
}

func (p *PriviligeRepository) FindPriviligeByName(name *string) *Privilige {
	var privilige Privilige
	p.db.Where("privilige=?", name).First(&privilige)
	return &privilige
}
