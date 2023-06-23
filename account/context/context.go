package context

import (
	"account/config"
	db "account/database"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	Db            *gorm.DB
	Config        *config.Config
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	err           error
	file          *os.File
)

func init() {


	Config = config.Load("application.yml")

	dbinfo := Config.Database
	Db, err = db.ConnectToDB(dbinfo.Username, dbinfo.Password, dbinfo.Dbname, dbinfo.Host, dbinfo.Port, dbinfo.Driver)
	if err != nil {
		log.Fatal(err)
	}

	file, err = os.OpenFile(Config.File.Logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("you can show logs in file: " + Config.File.Logpath)
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

}
