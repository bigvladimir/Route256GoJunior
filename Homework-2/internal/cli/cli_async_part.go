package cli

import (
	"Homework-2/internal/model"
	"Homework-2/internal/pvz_storage"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type access struct {
	commonWG     sync.WaitGroup
	mxForOldCode sync.Mutex
}

func UserInput() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	checkInputErr := make(chan struct{}, 1)

	PVZstor := pvz_storage.NewPvzStorage(pvzStorageFile)

	var control access
	defer func() {
		fmt.Println("Ожидание завершения всех процессов...")
		control.commonWG.Wait()
		PVZstor.Wg.Wait()
	}()

	// не менял структуру ввода чтобы не переделывать контроллер прошлого дз
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Программа запущена. Для получения справки напишите:", HelpCommandText, '\n')
		for {
			scanner.Scan()
			if scanner.Err() != nil {
				checkInputErr <- struct{}{}
				break
			}
			parts := strings.Fields(scanner.Text())

			switch len(parts) {
			case 0:
				continue
			case 1:
				switch parts[0] {
				case HelpCommandText:
					fmt.Println(helpText)
				default:
					fmt.Println("Неизвестная команда, используйте:", HelpCommandText)
				}
			default:
				id, err := strconv.Atoi(parts[0])
				if err != nil {
					fmt.Println("Не удалось привести строковую запись id в число:", err)
					break
				}
				if id <= 0 {
					fmt.Println("Некорректный id.")
					break
				}

				switch parts[1] {
				case checkPvzCommandText:
					control.commonWG.Add(1)
					go func() {
						defer control.commonWG.Done()

						if len(parts) != 2 {
							fmt.Println("Некорректное количество аргументов при просмотре информации о ПВЗ, можете открыть", HelpCommandText)
							return
						}

						PVZstor.AddWG()
						p, err := PVZstor.Read(id)
						if err != nil {
							fmt.Println("Ошибка при чтении базы данных ПВЗ:", err)
							return
						}
						if p == (model.Pvz{}) {
							fmt.Printf("ПВЗ id %d не найден в базе.\n", id)
							return
						}
						fmt.Printf("ПВЗ id %d Название: %s Адресс: %s Контакты: %s\n", p.PvzID, p.Name, p.Adress, p.Contacts)
					}()
				case newPvzCommandText:
					control.commonWG.Add(1)
					go func() {
						defer control.commonWG.Done()

						if len(parts) < 7 {
							fmt.Println("Слишком мало аргументов при добавлении ПВЗ, можете открыть", HelpCommandText)
							return
						}
						var pvzInput [3]string
						j := 0
						for i := 2; i < len(parts) && j < 3; i++ {
							if parts[i] == "/" {
								j++
								continue
							}
							if len(pvzInput[j]) > 0 {
								pvzInput[j] += " "
							}
							pvzInput[j] += parts[i]
						}
						if j > 2 && j < 2 {
							fmt.Println("Некорректное разделение аргументов при добавлении ПВЗ.", HelpCommandText)
							return
						}
						if len(pvzInput[0]) == 0 || len(pvzInput[1]) == 0 || len(pvzInput[2]) == 0 {
							fmt.Println("Неполные данные при добавлении ПВЗ.", HelpCommandText)
							return
						}
						PVZstor.AddWG()
						err := PVZstor.Write(model.Pvz{
							PvzID:    id,
							Name:     pvzInput[0],
							Adress:   pvzInput[1],
							Contacts: pvzInput[2],
						})
						if err != nil {
							fmt.Println("Ошибка при записи в базу данных ПВЗ:", err)
							return
						}
						fmt.Printf("ПВЗ id %d успешно добавлен в базу.\n", id)
					}()
				default:
					// старый код использует другую бд, поэтому он может запускаться
					// независимо от чтения/записи ПВЗ, разве что ему нужна операция Read()
					// чтобы проверить наличие пвз, но она сама залочится,
					// но самому старому коду нужен отдельный мьютекс,
					// не совсем аккуратно но не хочется прошлую домашку переписывать
					control.commonWG.Add(1)
					go func() {
						defer control.commonWG.Done()

						PVZstor.AddWG()
						p, err := PVZstor.Read(id)
						if err != nil {
							fmt.Println("Ошибка при проверке наличия ПВЗ в базе:", err)
							return
						}
						if p.PvzID <= 0 {
							fmt.Println("Невозможно выполнить запрос к базе данных заказов, если не существует записи о ПВЗ.")
							return
						}

						control.mxForOldCode.Lock()
						defer control.mxForOldCode.Unlock()

						err = orderRequest(parts)
						if err != nil {
							fmt.Println("Ошибка при работе с базой данных заказов:", err)
							return
						}
					}()
				}
			}
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

	return nil
}
