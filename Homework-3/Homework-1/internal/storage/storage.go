package storage

import (
	"Homework-1/internal/model"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

const storageFile = "../data/data.json"

var ZeroTime = time.Time{}

type Storage struct {
	storage *os.File
}

// New инициализирует Storage
func New() (Storage, error) {
	file, err := os.OpenFile(storageFile, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}

	return Storage{storage: file}, nil
}

// Close() закрывает файл
func (s Storage) Close() error {
	return s.storage.Close()
}

// Add добавляет новую запись, возвращает ошибку есои запись существует
func (s Storage) Add(input model.OrderInput) error {
	checkUnique, err := s.Get("orderID", input.OrderID)
	if err != nil {
		return err
	}
	if len(checkUnique) > 0 {
		return errors.New("Заказ уже записан в базу.")
	}

	all, err := s.getAll()
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// Get возвращает слайс заказов по фильтру
func (s Storage) Get(filter string, id ...int) ([]model.Order, error) {
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
	return orders, err
}

// Update обновляет запись в базе, если она присутствовала, иначе возвращает ошибку
func (s Storage) Update(input model.Order) error {
	checkUnique, err := s.Get("orderID", input.OrderID)
	if err != nil {
		return err
	}
	if len(checkUnique) == 0 {
		return errors.New("Попытка перезаписать несуществующую запись.")
	}

	all, err := s.getAll()
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// Delete удаляет запись из базы, если она присутствовала, иначе возвращает ошибку
func (s Storage) Delete(id int) error {
	checkUnique, err := s.Get("orderID", id)
	if err != nil {
		return err
	}
	if len(checkUnique) == 0 {
		return errors.New("Попытка удалить несуществующую запись.")
	}

	all, err := s.getAll()
	if err != nil {
		return err
	}
	for i, order := range all {
		if order.OrderID == id {
			all = append(all[:i], all[i+1:]...)
			break
		}
	}

	err = s.overwriteBytes(all)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) getByFilter(f func(orderDTO) bool) ([]model.Order, error) {
	all, err := s.getAll()
	if err != nil {
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

func (s Storage) getAll() ([]orderDTO, error) {
	_, err := s.storage.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(s.storage)
	bytesReader, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	orders := make([]orderDTO, 0)
	if len(bytesReader) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(bytesReader, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s Storage) overwriteBytes(orders []orderDTO) error {
	bytesReader, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	_, err = s.storage.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = s.storage.Truncate(0)
	if err != nil {
		return err
	}
	err = os.WriteFile(storageFile, bytesReader, 0777)
	if err != nil {
		return err
	}

	return nil
}
