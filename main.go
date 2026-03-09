package main

import (
	"cli_todo/httpServer"
	"cli_todo/storage/postgres"
	"log"
)

func main() {

	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewTaskRepository(db)
	server := httpServer.New(repo)
	err = server.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}