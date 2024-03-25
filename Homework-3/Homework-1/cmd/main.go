package main

import (
	"Homework-1/internal/cli"
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"fmt"
)

func main() {
	stor, err := storage.New()
	if err != nil {
		fmt.Println("Не удалось подключиться к хранилищу.")
		return
	}
	defer stor.Close()

	serv := service.New(&stor)

	err = cli.ExecuteRequestFromArgs(serv)
	if err != nil {
		fmt.Println("Во время выполнения запроса произошла ошибка:")
		fmt.Println(err)
		fmt.Println("Для получения справки используйте help.")
	} else {
		fmt.Println("Запрос выполнен успешно.")
	}
}
