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
	env := "local"
	SERVER := os.Getenv("DATABASE_SERVER")
	USERID := os.Getenv("USERID")
	DATABASE := os.Getenv("DATABASE")
	PASSWORD := os.Getenv("PASSWORD")
	PORT := os.Getenv("DB_PORT")

	var dbStatus obj.ConfigStruct
	var connectionString string

	if env == "live" {
		fmt.Println("connected live")
		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", USERID, PASSWORD, SERVER, "25060", DATABASE)
	} else {
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", SERVER, USERID, PASSWORD, DATABASE, PORT)
	}

	logrus.Info(connectionString)
	DbGorm, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, NoLowerCase: true,
	}, PrepareStmt: false})

	if err != nil {
		fmt.Sprintln(err.Error())
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
		DbGorm.AutoMigrate(&dbSchema.TblRole{})
		DbGorm.AutoMigrate(&dbSchema.TblGroupType{})
		DbGorm.AutoMigrate(&dbSchema.TblGroupUser{})
		DbGorm.AutoMigrate(&dbSchema.TblOrderItem{})
		DbGorm.AutoMigrate(&dbSchema.TblStatus{})
		DbGorm.AutoMigrate(&dbSchema.TblCart{})
		DbGorm.AutoMigrate(&dbSchema.TblTitle{})
		DbGorm.AutoMigrate(&dbSchema.TblCartItem{})
		DbGorm.AutoMigrate(&dbSchema.TblCartMember{})
		DbGorm.AutoMigrate(&dbSchema.TblProduct{})
		DbGorm.AutoMigrate(&dbSchema.TblUser{})
		DbGorm.AutoMigrate(&dbSchema.TblClient{})
		DbGorm.AutoMigrate(&dbSchema.TblCartType{})
	}

	dbStatus.IsDropExistingTables = false
	dbStatus.CreateTable = false
	domarchal, _ := json.Marshal(dbStatus)
	_ = os.WriteFile("config.json", domarchal, 0400)
	return DbGorm
}
