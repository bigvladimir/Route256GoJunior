package cli

import (
	"errors"
	"fmt"
	"homework/internal/interactive/model"
	"homework/internal/interactive/order_storage"
	"homework/internal/interactive/service"
	"strconv"
	"time"
)

// orderRequest вызывает команды service
func orderRequest(args []string) error {
	if len(args) == 0 {
		return errors.New("Не введено id ПВЗ и(или) команда.")
	}

	if len(args) == 1 {
		return errors.New("Не введена команда.")
	}

	pvzID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
	}
	if pvzID <= 0 {
		return errors.New("Некорректный id ПВЗ.")
	}
	stor, err := order_storage.New(pvzID)
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к базе данных заказов: %w", err)
	}
	defer stor.Close()

	serv := service.New(&stor)

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
			return fmt.Errorf("Не удалось принять заказ %d у курьера: %w", orderID, err)
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
		fmt.Printf("Заказ %d возвращён курьеру.\n", orderID)

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
		fmt.Println("Выдача обработана, заказы выданы:", orderID)

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
		fmt.Printf("Возврат заказа %d принят.\n", orderID)

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

	default:
		return errors.New("Такой команды нет.")
	}

	return nil
}

// TODO будет перемешиваться при асинхронном выводе
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
