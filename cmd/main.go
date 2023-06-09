package main

import (
	"fmt"
	"os"

	"github.com/itsshashank/identity-reconciliation/api"
	"github.com/itsshashank/identity-reconciliation/db"
)

func main() {
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	userStore := db.NewPostgresUserStore(dsn)

	server := api.NewServer(userStore)
	server.Listen(listenAddr)
}
