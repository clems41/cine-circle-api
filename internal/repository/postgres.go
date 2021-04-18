package repository

import (
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = "5432"
	defaultUser = "postgres"
	defaultDbName = "cine-circle"
	defaultPassword = "postgres"
	defaultDebug = "true"
	defaultDetailedLogs = "true"
)

const (
	EnvHost = "DB_HOST"
	EnvPort = "DB_PORT"
	EnvUser = "DB_USER"
	EnvDbName = "DB_NAME"
	EnvPassword = "DB_PASSWORD"
	EnvDebug = "DB_DEBUG"
	EnvDetailedLogs = "DB_LOG"
)

type Database struct {
	db *gorm.DB
}


type PostgresConfig struct {
	Host            string
	User            string
	Password        string
	Port            string
	DbName          string
	ExtraConfigs    string
	Debug           bool
	DetailedLogs    bool
	ApplicationName string
}

func (postgresConfig PostgresConfig) DataSourceName() string {
	dataSourceName := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v application_name='%v' %v",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.DbName,
		postgresConfig.ApplicationName,
		postgresConfig.ExtraConfigs)

	return dataSourceName
}



func OpenConnection() (db *gorm.DB, err error) {
	pgConfig := PostgresConfig{
		Host:            utils.GetDefaultOrFromEnv(defaultHost, EnvHost),
		User:            utils.GetDefaultOrFromEnv(defaultUser, EnvUser),
		Password:        utils.GetDefaultOrFromEnv(defaultPassword, EnvPassword),
		Port:            utils.GetDefaultOrFromEnv(defaultPort, EnvPort),
		DbName:          utils.GetDefaultOrFromEnv(defaultDbName, EnvDbName),
		Debug:           utils.GetDefaultOrFromEnv(defaultDebug, EnvDebug) == "true",
		DetailedLogs:    utils.GetDefaultOrFromEnv(defaultDetailedLogs, EnvDetailedLogs) == "true",
		ExtraConfigs: 	 "sslmode=disable TimeZone=Pacific/Noumea",
		ApplicationName: "cine-circle-import",
	}
	gormCfg := gorm.Config{}
	if pgConfig.DetailedLogs {
		newLogger := gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      gormLogger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		)
		gormCfg.Logger = newLogger
	}
	db, err = gorm.Open(postgres.Open(pgConfig.DataSourceName()), &gormCfg)
	if err != nil {
		return
	}

	if db == nil {
		return nil, typedErrors.ErrRepositoryIsNil
	}

	if pgConfig.Debug {
		db = db.Debug()
	}
	return
}

func CloseConnection(DB *gorm.DB) (err error) {
	db, err := DB.DB()
	if err != nil {
		return
	}
	return db.Close()
}