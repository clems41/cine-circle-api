package db

import (
	"cine-circle/internal/model"
	"cine-circle/internal/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
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

func OpenConnection() (*Database, model.CustomError) {
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
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		)
		gormCfg.Logger = newLogger
	}
	database, err := gorm.Open(postgres.Open(pgConfig.DataSourceName()), &gormCfg)
	if err != nil {
		panic(err)
	}

	if database == nil {
		panic(model.ErrInternalDatabaseIsNil)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	//database.DB().SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//database.DB().SetMaxOpenConns(10)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//database.DB().SetConnMaxLifetime(time.Hour)
	return &Database{db: database}, model.NewCustomError(err, http.StatusInternalServerError)
}

func (db *Database) Close() model.CustomError {
	if db.db == nil {
		return model.NewCustomError(nil, http.StatusInternalServerError)
	}
	sqlDB, err := db.db.DB()
	if err != nil {
		return model.NewCustomError(err, http.StatusInternalServerError)
	}
	return model.NewCustomError(sqlDB.Close(), http.StatusInternalServerError)
}