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
	case "takeorder":
		if len(args) != 5 {
			return errors.New("Некорректное количество аргументов.")
		}
		orderID, err := atoiCustom(args[2])
		if err != nil {
			return err
		}
		customerID, err := atoiCustom(args[3])
		if err != nil {
			return err
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
			return err
		}
		fmt.Println("Заказ принят.")

	case "returnorder":
		if len(args) != 3 {
			return errors.New("Некорректное количество аргументов.")
		}
		orderID, err := atoiCustom(args[2])
		if err != nil {
			return err
		}

		err = serv.ReturnOrderToCourier(orderID)
		if err != nil {
			return err
		}
		fmt.Println("Заказ возвращён курьеру.")

	case "giveorder":
		if len(args) < 3 {
			return errors.New("Некорректное количество аргументов.")
		}
		var orderID []int
		for i := 2; i < len(args); i++ {
			id, err := atoiCustom(args[i])
			if err != nil {
				return err
			}
			orderID = append(orderID, id)
		}

		err := serv.GiveOrderToCustomer(orderID)
		if err != nil {
			return err
		}
		fmt.Println("Выдача обработана.")

	case "orderlist":
		if len(args) < 3 || len(args) > 6 {
			return errors.New("Некорректное количество аргументов.")
		}
		limit := 0
		isInStock := false
		customerID, err := atoiCustom(args[2])
		if err != nil {
			return err
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
			limit, err = atoiCustom(args[4])
			if err != nil {
				return err
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
			limit, err = atoiCustom(args[limitPos+1])
			if err != nil {
				return err
			}
		}

		orders, err := serv.GetCustomerOrderList(customerID, limit, isInStock)
		if err != nil {
			return err
		}
		for _, o := range orders {
			printOrder(o)
		}

	case "takerefund":
		if len(args) != 4 {
			return errors.New("Некорректное количество аргументов.")
		}
		customerID, err := atoiCustom(args[2])
		if err != nil {
			return err
		}
		orderID, err := atoiCustom(args[3])
		if err != nil {
			return err
		}

		err = serv.TakeRefundFromCustomer(customerID, orderID)
		if err != nil {
			return err
		}
		fmt.Println("Возврат принят.")

	case "refundlist":
		if len(args) != 4 {
			return errors.New("Некорректное количество аргументов.")
		}
		pageNum, err := atoiCustom(args[2])
		if err != nil {
			return err
		}
		pageSize, err := atoiCustom(args[3])
		if err != nil {
			return err
		}

		orders, err := serv.GetRefundList(pageNum, pageSize)
		if err != nil {
			return err
		}
		for _, o := range orders {
			printOrder(o)
		}

	case "help":
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

func atoiCustom(s string) (int, error) {
	answ, err := strconv.Atoi(s)
	if err != nil {
		err := errors.New("Не удалось преобразовать аргумент в число.")
		return -1, err
	}
	return answ, nil
}
