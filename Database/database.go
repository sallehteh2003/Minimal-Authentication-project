package Database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/Config"
)

type DB struct {
	cfg Config.Config
	sql *gorm.DB
}

func CreateAndConnectToDb(cfg Config.Config) (*DB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	// Create a new connection
	db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{
		cfg: cfg,
		sql: db,
	}, nil
}

func (gdb *DB) CreateModel() error {

	err := gdb.sql.AutoMigrate(User{})
	if err != nil {
		return err
	}

	return nil

}
