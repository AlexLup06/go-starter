package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Config is the configuration for the database.
type Config struct {
	Host     string `env:"POSTGRESQL_HOST, required"`
	Port     int    `env:"POSTGRESQL_PORT, required"`
	Username string `env:"POSTGRESQL_USERNAME, required"`
	Password string `env:"POSTGRESQL_PASSWORD, required"`
	Database string `env:"POSTGRESQL_DATABASE, required"`
	Timezone string `env:"POSTGRESQL_TIMEZONE, required"`
	Schema   string `env:"POSTGRESQL_SCHEMA, required"`
	LogSql   bool   `env:"POSTGRESQL_LOG_SQL, required"`
}

type Client struct {
	db *gorm.DB
}

func NewClient(cfg Config) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s search_path=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.Timezone, cfg.Schema)

	var lg gormlogger.Interface
	if cfg.LogSql {
		lg = gormlogger.Default.LogMode(gormlogger.Info)
	}

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, fmt.Errorf("can't load location %s: %w", cfg.Timezone, err)
	}

	gormConfig := gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
		Logger:                 lg,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Schema + ".",
			SingularTable: false,
		},
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	c := Client{db: gormDB}

	return &c, nil
}

func (c *Client) GormDb() *gorm.DB {
	return c.db
}
