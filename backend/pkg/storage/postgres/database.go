package postgres

import (
	"context"
	"fmt"
	"github.com/chempik1234/availability-checker-web/config"
	"github.com/jackc/pgx/v5"
	"log"
)

// DBInstance is a structure for storing a database definition
type DBInstance struct {
	Db pgx.Conn
}

// NewDBInstance try to connect to the database and create a new DBInstance to use in repos
func NewDBInstance(ctx context.Context, config *config.DB) (*DBInstance, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		config.DbHost,
		config.DbUser,
		config.DbName,
		config.DbPassword,
	)

	db, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		log.Fatal("Failed to connect to databse. \n", err)
	}

	log.Println("connected to the postgresql database")

	DB := DBInstance{*db}
	return &DB, nil
}
