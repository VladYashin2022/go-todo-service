package main

import (
	"cli_todo/cli"
	"cli_todo/httpServer"
	"cli_todo/service"
	"cli_todo/storage"
	"cli_todo/storage/postgres"
	"flag"
	"fmt"
	"log"
)

func main() {
	WriteAllTaskJson()

	service.TasksId = service.FindMaxID(service.AllTasks)

	var httpMode = flag.Bool("http", false, "run http server with pgsql")
	flag.Parse()

	// запуск http
	if *httpMode {
		db, err := postgres.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		repo := postgres.NewTaskRepository(db)
		server := httpServer.New(repo)
		err = server.Run("localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := cli.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// в срез AllTasks вносим все Task из storage.json, если такой файл имеется
func WriteAllTaskJson() {
	var err error
	service.AllTasks, err = storage.AllTasksWriter()
	if err != nil {
		fmt.Println(err)
	}
}
