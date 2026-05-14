package db

import (
	"context"
	"fmt"
	"logarda/internal/config"
	"logarda/internal/model"

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

func VerifyUserCredentials(ctx context.Context, username, hashedPassword string, user *model.User) error {
	err := DB.QueryRow(ctx, "SELECT username, password FROM users WHERE username=$1 AND password=$2;",
		username, hashedPassword).Scan(&user.Username, &user.Password)
	return err
}

func SaveAWSCredentials(ctx context.Context, encryptedID, encryptedSecret, region, username string) error {
	_, err := DB.Exec(ctx, "UPDATE users SET awskeyid = $1, awskeysecret=$2, awsregion=$3 WHERE username=$4;",
		encryptedID, encryptedSecret, region, username)
	return err
}
