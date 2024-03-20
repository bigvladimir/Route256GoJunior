package order_storage

import (
	"Homework-2/internal/model"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// нужно дописать id и расширение
const storageFileSuff = "../data/orders/pvz"

var ZeroTime = time.Time{}

type OrderStorage struct {
	storage  *os.File
	filePath string
}

// New инициализирует Storage
// информация об айди пвз проверяется на более верхних уровнях,
func New(pvzID int) (OrderStorage, error) {
	filePath := storageFileSuff + strconv.Itoa(pvzID) + ".json"
	file, err := os.OpenFile(filePath, os.O_CREATE, 0777)
	if err != nil {
		err = fmt.Errorf("Не удалось открыть файл: %w", err)
		return OrderStorage{}, err
	}

	return OrderStorage{storage: file, filePath: filePath}, nil
}

// Close() закрывает файл
func (s OrderStorage) Close() error {
	return s.storage.Close()
}

// Add добавляет новую запись, возвращает ошибку есои запись существует
func (s OrderStorage) Add(input model.OrderInput) error {
	checkUnique, err := s.Get("orderID", input.OrderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить записи по ID заказа: %w", err)
	}
	if len(checkUnique) > 0 {
		return errors.New("Заказ уже записан в базу.")
	}

	all, err := s.getAll()
	if err != nil {
		return fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
	}

	newOrder := orderDTO{
		OrderID:         input.OrderID,
		CustomerID:      input.CustomerID,
		StorageLastTime: input.StorageLastTime,
		IsCompleted:     false,
		CompleteTime:    ZeroTime,
		IsRefunded:      false,
		ArrivalTime:     time.Now(),
	}

	all = append(all, newOrder)

	err = s.overwriteBytes(all)
	if err != nil {
		return fmt.Errorf("Не удалось записать данные в базу данных: %w", err)
	}

	return nil
}

// Get возвращает слайс заказов по фильтру
func (s OrderStorage) Get(filter string, id ...int) ([]model.Order, error) {
	var filterFunc func(orderDTO) bool
	switch filter {
	case "orderID":
		if len(id) != 1 {
			err := errors.New("Неправильное количество аргументов для фильтра orderID в Get().")
			return nil, err
		}
		filterFunc = func(o orderDTO) bool { return o.OrderID == id[0] }
	case "customerID":
		if len(id) != 1 {
			err := errors.New("Неправильное количество аргументов для фильтра customerID в Get().")
			return nil, err
		}
		filterFunc = func(o orderDTO) bool { return o.CustomerID == id[0] }
	case "isRefunded":
		if len(id) != 0 {
			err := errors.New("Неправильное количество аргументов для фильтра isRefunded в Get().")
			return nil, err
		}
		filterFunc = func(o orderDTO) bool { return o.IsRefunded == true }
	default:
		err := errors.New("Непредусмотренный фильтр в Get().")
		return nil, err
	}

	orders, err := s.getByFilter(filterFunc)
	if err != nil {
		err = fmt.Errorf("Не удалось получить данные по фильтру: %w", err)
		return nil, err
	}

	return orders, nil
}

// Update обновляет запись в базе, если она присутствовала, иначе возвращает ошибку
func (s OrderStorage) Update(input model.Order) error {
	checkUnique, err := s.Get("orderID", input.OrderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить записи по ID заказа: %w", err)
	}
	if len(checkUnique) == 0 {
		return errors.New("Попытка перезаписать несуществующую запись.")
	}

	all, err := s.getAll()
	if err != nil {
		return fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
	}
	for i, order := range all {
		if order.OrderID == input.OrderID {
			all[i] = orderDTO{
				OrderID:         input.OrderID,
				CustomerID:      input.CustomerID,
				StorageLastTime: input.StorageLastTime,
				IsCompleted:     input.IsCompleted,
				CompleteTime:    input.CompleteTime,
				IsRefunded:      input.IsRefunded,
				ArrivalTime:     input.ArrivalTime,
			}
			break
		}
	}

	err = s.overwriteBytes(all)
	if err != nil {
		return fmt.Errorf("Не удалось записать данные в базу данных: %w", err)
	}

	return nil
}

// Delete удаляет запись из базы, если она присутствовала, иначе возвращает ошибку
func (s OrderStorage) Delete(id int) error {
	checkUnique, err := s.Get("orderID", id)
	if err != nil {
		return fmt.Errorf("Не удалось получить записи по ID заказа: %w", err)
	}
	if len(checkUnique) == 0 {
		return errors.New("Попытка удалить несуществующую запись.")
	}

	all, err := s.getAll()
	if err != nil {
		return fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
	}
	for i, order := range all {
		if order.OrderID == id {
			all = append(all[:i], all[i+1:]...)
			break
		}
	}

	err = s.overwriteBytes(all)
	if err != nil {
		return fmt.Errorf("Не удалось записать данные в базу данных: %w", err)
	}

	return nil
}

func (s OrderStorage) getByFilter(f func(orderDTO) bool) ([]model.Order, error) {
	if f == nil {
		err := errors.New("Передана пустая функция фильтрации в getByFilter().")
		return nil, err
	}

	all, err := s.getAll()
	if err != nil {
		err = fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
		return nil, err
	}

	filtered := make([]model.Order, 0)
	for _, order := range all {
		if f(order) {
			filtered = append(filtered, model.Order{
				OrderID:         order.OrderID,
				CustomerID:      order.CustomerID,
				StorageLastTime: order.StorageLastTime,
				IsCompleted:     order.IsCompleted,
				CompleteTime:    order.CompleteTime,
				IsRefunded:      order.IsRefunded,
				ArrivalTime:     order.ArrivalTime,
			})
		}
	}

	return filtered, nil
}

func (s OrderStorage) getAll() ([]orderDTO, error) {
	_, err := s.storage.Seek(0, io.SeekStart)
	if err != nil {
		err = fmt.Errorf("Ошибка Seek() при перемещении на начальную позицию в файле: %w", err)
		return nil, err
	}
	reader := bufio.NewReader(s.storage)
	bytesReader, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("Ошибка при вызове ReadAll(): %w", err)
		return nil, err
	}
	orders := make([]orderDTO, 0)
	if len(bytesReader) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(bytesReader, &orders)
	if err != nil {
		err = fmt.Errorf("Ошибка при распаковке json (json.Unmarshal()): %w", err)
		return nil, err
	}

	return orders, nil
}

func (s OrderStorage) overwriteBytes(orders []orderDTO) error {
	bytesReader, err := json.Marshal(orders)
	if err != nil {
		return fmt.Errorf("Ошибка при приведении данных в json (json.Marshal()): %w", err)
	}

	err = os.WriteFile(s.filePath, bytesReader, 0777)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове WriteFile(): %w", err)
	}

	return nil
}
