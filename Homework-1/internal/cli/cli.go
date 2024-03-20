package cli

import (
	"Homework-1/internal/model"
	"Homework-1/internal/service"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ExecuteRequestFromArgs по информации из аргументов утилиты вызывает команды service
func ExecuteRequestFromArgs(serv service.Service) error {
	args := os.Args

	if len(args) == 1 {
		return errors.New("Не введена команда.")
	}

	switch args[1] {
	case takeOrderCommandText:
		if len(args) != 5 {
			return errors.New("Некорректное количество аргументов.")
		}
		orderID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		customerID, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		layout := "2006-01-02"
		lastDate, err := time.Parse(layout, args[4])
		if err != nil {
			return errors.New("Некорректная дата.")
		}

		order := model.OrderInput{
			OrderID:         orderID,
			CustomerID:      customerID,
			StorageLastTime: lastDate,
		}
		err = serv.TakeOrderFromCourier(order)
		if err != nil {
			return fmt.Errorf("Не удалось принять заказ у курьера: %w", err)
		}
		fmt.Println("Заказ принят.")

	case returnOrderCommandText:
		if len(args) != 3 {
			return errors.New("Некорректное количество аргументов.")
		}
		orderID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		err = serv.ReturnOrderToCourier(orderID)
		if err != nil {
			return fmt.Errorf("Не удалось вернуть заказ курьеру: %w", err)
		}
		fmt.Println("Заказ возвращён курьеру.")

	case giveOrderCommandText:
		if len(args) < 3 {
			return errors.New("Некорректное количество аргументов.")
		}
		var orderID []int
		for i := 2; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
			orderID = append(orderID, id)
		}

		err := serv.GiveOrderToCustomer(orderID)
		if err != nil {
			return fmt.Errorf("Не удалось выдать заказы клиенту: %w", err)
		}
		fmt.Println("Выдача обработана.")

	case orderListCommandText:
		if len(args) < 3 || len(args) > 6 {
			return errors.New("Некорректное количество аргументов.")
		}
		limit := 0
		isInStock := false
		customerID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		if len(args) == 4 {
			if args[3] != "instock" {
				return errors.New("Некорректные дополнительные аргументы.")
			}
			isInStock = true
		}
		if len(args) == 5 {
			if args[3] != "limit" {
				return errors.New("Некорректные дополнительные аргументы.")
			}
			limit, err = strconv.Atoi(args[4])
			if err != nil {
				return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
		}
		if len(args) == 6 {
			limitPos := 3
			if args[3] == "instock" {
				limitPos = 4
			} else if args[5] != "instock" {
				return errors.New("Некорректные дополнительные аргументы.")
			}
			isInStock = true
			if args[limitPos] != "limit" {
				return errors.New("Некорректные дополнительные аргументы.")
			}
			limit, err = strconv.Atoi(args[limitPos+1])
			if err != nil {
				return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
		}

		orders, err := serv.GetCustomerOrderList(customerID, limit, isInStock)
		if err != nil {
			return fmt.Errorf("Не удалось получить список заказов: %w", err)
		}
		for _, o := range orders {
			printOrder(o)
		}

	case takeRefundCommandText:
		if len(args) != 4 {
			return errors.New("Некорректное количество аргументов.")
		}
		customerID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		orderID, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		err = serv.TakeRefundFromCustomer(customerID, orderID)
		if err != nil {
			return fmt.Errorf("Не удалось принять возврат от клиента: %w", err)
		}
		fmt.Println("Возврат принят.")

	case refundListCommandText:
		if len(args) != 4 {
			return errors.New("Некорректное количество аргументов.")
		}
		pageNum, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		pageSize, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		orders, err := serv.GetRefundList(pageNum, pageSize)
		if err != nil {
			return fmt.Errorf("Не удалось получить список возвратов: %w", err)
		}
		for _, o := range orders {
			printOrder(o)
		}

	case HelpCommandText:
		fmt.Println(helpText)

	default:
		return errors.New("Такой команды нет.")
	}

	return nil
}

func printOrder(order model.Order) {
	fmt.Println("id:", order.OrderID,
		"; customer:", order.CustomerID,
		"; arrived:", order.ArrivalTime.Year(),
		order.ArrivalTime.Month(), order.ArrivalTime.Day(),
		"; stored until:", order.StorageLastTime.Year(),
		order.StorageLastTime.Month(), order.StorageLastTime.Day(),
		"; completed:", order.IsCompleted,
		"; refunded:", order.IsRefunded)
}
