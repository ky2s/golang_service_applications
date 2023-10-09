package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configurations exported
type Configurations struct {
	Server          ServerConfigurations
	Database        DatabaseConfigurations
	EXAMPLE_PATH    string
	EXAMPLE_VAR     string
	SSL_PRIVATE_KEY string
	SSL_PUBLIC_KEY  string

	ENV_TYPE string

	OSS_ACCESS_KEY_ID     string
	OSS_SECRET_ACCESS_KEY string
	OSS_REGION            string
	OSS_BUCKET            string
	OSS_ENDPOINT          string
	OSS_URL               string

	LINK_ONESINYAL   string
	APP_ID_ONESINYAL string
	KEY_ONESINYAL    string

	SMTP_HOST   string
	SMTP_PORT   int
	EMAIL       string
	SENDER_NAME string
	PASSWORD    string

	LINKVERIFY string

	SHORTEN_BASE_URL string
	SHORTEN_KEY      string

	LINK_EXTERNAL string
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Hostname string
	Port     int
	Ssl_Port int
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func ConnectDB(config Configurations) *gorm.DB {

	URL := "host=" + config.Database.DBHost + " user=" + config.Database.DBUser + " password=" + config.Database.DBPassword + " dbname=" + config.Database.DBName + " port=" + config.Database.DBPort + " sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(URL), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
