package db

import (
	"context"
	"fmt"
	"logarda/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var err error
var NoRows error

func CreatePostgreSQLPool() {
	ctx := context.Background()
	NoRows = pgx.ErrNoRows

	DB, err = pgxpool.New(ctx, config.POSTGRES_URL) // assigns pool to global variable
	if err != nil {
		DB = nil
		fmt.Println("Error connecting to Postgres database!")
		return
	}
}
