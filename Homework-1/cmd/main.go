package main

import (
	"Homework-1/internal/cli"
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"fmt"
	"log"
)

func main() {
	stor, err := storage.New()
	if err != nil {
		log.Fatal("Не удалось подключиться к хранилищу:", err)
	}
	defer stor.Close()

	serv := service.New(&stor)

	err = cli.ExecuteRequestFromArgs(serv)
	if err != nil {
		fmt.Printf("Во время выполнения запроса произошла ошибка:\n%v\n", err)
		fmt.Println("Для получения справки используйте:", cli.HelpCommandText)
	} else {
		fmt.Println("Запрос выполнен успешно.")
	}
}
