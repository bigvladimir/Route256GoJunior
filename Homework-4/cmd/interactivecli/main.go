package main

import (
	"bufio"
	"context"
	"fmt"
	"homework/internal/app/clihandler"
	"homework/internal/app/core"
	"homework/internal/app/orders"
	"homework/internal/app/pvz"
	"homework/internal/pkg/db"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type access struct {
	commonWG sync.WaitGroup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}
	defer database.GetPool(ctx).Close()

	packVariants := map[string]orders.PackageVariant{
		orders.BagVariantName:  orders.NewBagPackage(),
		orders.BoxVariantName:  orders.NewBoxPackage(),
		orders.FilmVariantName: orders.NewFilmPackage(),
	}
	orderService := orders.NewService(orders.NewOrderStorage(database), packVariants)
	pvzService := pvz.NewService(pvz.NewPvzStorage(database))

	service := core.NewCoreService(orderService, pvzService)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	checkInputErr := make(chan struct{}, 1)

	var control access
	defer func() {
		fmt.Println("Ожидание завершения всех процессов...")
		control.commonWG.Wait()
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Программа запущена. Для получения справки напишите:", clihandler.HelpCommandText)
		for {
			scanner.Scan()
			if scanner.Err() != nil {
				checkInputErr <- struct{}{}
				break
			}
			text := scanner.Text()
			control.commonWG.Add(1)
			go func(command string) {
				defer control.commonWG.Done()
				// игнорируем пустой ввод, чтобы не засорять консоль
				if len(command) == 0 {
					return
				}
				answer, err := clihandler.ExecCommand(ctx, service, command)
				if err != nil {
					fmt.Println("Во время выполнения запроса произошла ошибка:", err)
					return
				}
				fmt.Printf("Результат выполнения запроса:\n%s\n", answer)
			}(text)
		}
	}()

	checkAllErr := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-sig:
				fmt.Println("Пришёл сигнал завершения работы.")
				checkAllErr <- struct{}{}
			case <-checkInputErr:
				fmt.Println("Ошибка функции ввода.")
				checkAllErr <- struct{}{}
			}
		}
	}()

	<-checkAllErr
	fmt.Println("Программа завершается...")
}
