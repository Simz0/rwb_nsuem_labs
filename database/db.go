package database

import (
	"log"
	"os"

	"labs/l0/database/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error. ENV file is undefined")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Error with database connect")
	}

	DB = db
	log.Println("Connect to DB is successfull")
}

func InitDatabase() error {
	err := DB.AutoMigrate(
		&models.Order{},
		&models.Delivery{},
		&models.Payment{},
		&models.Item{},
	)

	if err != nil {
		return err
	}

	return nil
}
