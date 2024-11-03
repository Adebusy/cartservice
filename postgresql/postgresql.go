package postgresql

import (
	"encoding/json"
	"fmt"
	"os"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/Adebusy/cartbackendsvc/obj"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DbGorm *gorm.DB
var err error

func GetDB() *gorm.DB {
	if loadEnv := godotenv.Load(); loadEnv != nil {
		ret := fmt.Sprintf("Unable to load environment variable. %s", loadEnv.Error())
		fmt.Println(ret)
	}
	env := "live"
	SERVER := os.Getenv("DATABASE_SERVER" + "_" + env)
	USERID := os.Getenv("USERID" + "_" + env)
	DATABASE := os.Getenv("DATABASE" + "_" + env)
	PASSWORD := os.Getenv("PASSWORD" + "_" + env)
	PORT := os.Getenv("PORT" + "_" + env)

	var dbStatus obj.ConfigStruct
	var connectionString string
	if env == "live" {
		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", USERID, PASSWORD, SERVER, PORT, DATABASE)
	} else {
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", SERVER, USERID, PASSWORD, DATABASE, PORT)
	}

	DbGorm, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, NoLowerCase: true,
	}})
	if err != nil {
		panic("failed to connect database")
	}

	read, err := os.ReadFile("config.json")
	if err != nil {
		logrus.Error(err)
	}
	if err := json.Unmarshal(read, &dbStatus); err != nil {
		logrus.Error(err)
	}

	if dbStatus.CreateTable {
		DbGorm.AutoMigrate(&dbSchema.TblStatus{})
		DbGorm.AutoMigrate(&dbSchema.TblCart{})
		DbGorm.AutoMigrate(&dbSchema.TblTitle{})
		DbGorm.AutoMigrate(&dbSchema.TblCartItem{})
		DbGorm.AutoMigrate(&dbSchema.TblCartMember{})
		DbGorm.AutoMigrate(&dbSchema.TblProduct{})
		DbGorm.AutoMigrate(&dbSchema.TblUser{})
		DbGorm.AutoMigrate(&dbSchema.TblCartType{})
	}
	dbStatus.IsDropExistingTables = false
	dbStatus.CreateTable = false
	domarchal, _ := json.Marshal(dbStatus)
	_ = os.WriteFile("config.json", domarchal, 0400)
	return DbGorm
}
