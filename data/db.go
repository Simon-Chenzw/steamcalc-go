package data

import (
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DEFAULT_DATABASE = "gorm.db"

var (
	once      sync.Once
	singleton DB
)

type DB struct {
	*gorm.DB
}

func GetDataBase() DB {
	once.Do(func() {
		slog.Info("Loading database...")
		viper.SetDefault("database", DEFAULT_DATABASE)

		db, err := gorm.Open(
			sqlite.Open(viper.GetString("database")),
			&gorm.Config{
				Logger: logger.Default.LogMode(logger.Error),
			},
		)
		if err != nil {
			panic(err)
		}

		singleton = DB{db}
		singleton.Setup()

		slog.Info("Database loaded.")
	})
	return singleton
}

func (db *DB) Setup() {
	slog.Info("Migrating database...")
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&MarketPrice{})

	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic(res.Error)
	}
}
