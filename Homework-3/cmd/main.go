package main

import (
	"homework/internal/interactive/cli"
	"homework/internal/sqlserver/setserver"
	"log"
	"os"
)

func main() {

	args := os.Args
	if len(args) != 2 {
		log.Fatal("Требуется один аргумент с вариантом работы (3, 2).")
	}

	switch args[1] {
	case "2":
		err := cli.UserInput()
		if err != nil {
			log.Fatal("Программа завершилась c ошибкой, причина:", err)
		}
	case "3":
		err := setserver.SetServer()
		if err != nil {
			log.Fatal("Сервер завершил работу c ошибкой, причина:", err)
		}
	default:
		log.Fatal("Допустимые аргументы: 3, 2.")
	}
}
