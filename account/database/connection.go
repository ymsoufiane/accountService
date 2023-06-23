package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(dbUser string, dbPassword string, dbName string, host string, port string, driver string) (*gorm.DB, error) {
	dbLogin := dbUser + ":" + dbPassword + "@tcp(" + host + ":" + port + ")/" + dbName //"soufiane:password@tcp(localhost:3306)/erp_medical"
	if driver == "mysql" {
		return gorm.Open(mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        dbLogin + "?charset=utf8&parseTime=True&loc=Local",
		}), &gorm.Config{})
	} else if driver == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require", host, port, dbUser, dbName, dbPassword)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	return nil,errors.New("driver value mysql or postgres")

}
