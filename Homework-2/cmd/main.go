package main

import (
	"Homework-2/internal/cli"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Программа запускается...")
	err := cli.UserInput()
	if err != nil {
		log.Fatal("Программа завершилась c ошибкой, причина:", err)
	}
	fmt.Println("Программа завершена.")
}
