package clihandler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	orders_dto "homework/internal/app/orders/dto"
	pvz_dto "homework/internal/app/pvz/dto"
	pvz_errors "homework/internal/app/pvz/errors"
)

type coreOps interface {
	TakeOrderFromCourier(ctx context.Context, order orders_dto.OrderInput) error
	ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error
	GiveOrderToCustomer(ctx context.Context, pvzID int, orderID []int) error
	TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error
	GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]orders_dto.Order, error)
	GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]orders_dto.Order, error)

	AddPvz(ctx context.Context, input pvz_dto.PvzInput) (int64, error)
	GetPvzByID(ctx context.Context, id int64) (pvz_dto.Pvz, error)
	UpdatePvz(ctx context.Context, input pvz_dto.Pvz) error
	DeletePvz(ctx context.Context, id int64) error
	ModifyPvz(ctx context.Context, input pvz_dto.Pvz) (int64, error)
}

// ExecCommand обрабатывает строковую команду и вызывает интерфейс корневого пакета
func ExecCommand(ctx context.Context, service coreOps, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 1 && parts[0] == HelpCommandText {
		return helpText, nil
	}

	if len(parts) < 2 {
		return "", errors.New("Неполная команда")
	}
	pvzID, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
	}
	if pvzID <= 0 {
		return "", errors.New("Некорректный id ПВЗ")
	}
	switch parts[1] {
	case takeOrderCommandText:
		if len(parts) != 7 && len(parts) != 8 {
			return "", errors.New("Некорректное количество аргументов")
		}
		orderID, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		customerID, err := strconv.Atoi(parts[3])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		layout := "2006-01-02"
		lastDate, err := time.Parse(layout, parts[4])
		if err != nil {
			return "", errors.New("Некорректная дата")
		}
		weight, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		price, err := strconv.Atoi(parts[6])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		pack := ""
		if len(parts) == 8 {
			pack = parts[7]
		}

		order := orders_dto.OrderInput{
			OrderID:         orderID,
			PvzID:           pvzID,
			CustomerID:      customerID,
			StorageLastTime: lastDate,
			PackageType:     pack,
			Weight:          weight,
			Price:           price,
		}
		err = service.TakeOrderFromCourier(ctx, order)
		if err != nil {
			return "", fmt.Errorf("Не удалось принять заказ %d у курьера: %w", orderID, err)
		}
		return "Заказ принят", nil
	case returnOrderCommandText:
		if len(parts) != 3 {
			return "", errors.New("Некорректное количество аргументов")
		}
		orderID, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		err = service.ReturnOrderToCourier(ctx, pvzID, orderID)
		if err != nil {
			return "", fmt.Errorf("Не удалось вернуть заказ курьеру: %w", err)
		}
		return fmt.Sprintf("Заказ %d возвращён курьеру", orderID), nil
	case giveOrderCommandText:
		if len(parts) < 3 {
			return "", errors.New("Некорректное количество аргументов")
		}
		var orderID []int
		for i := 2; i < len(parts); i++ {
			id, err := strconv.Atoi(parts[i])
			if err != nil {
				return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
			orderID = append(orderID, id)
		}

		err := service.GiveOrderToCustomer(ctx, pvzID, orderID)
		if err != nil {
			return "", fmt.Errorf("Не удалось выдать заказы клиенту: %w", err)
		}
		return fmt.Sprint("Выдача обработана, заказы выданы:", orderID), nil
	case orderListCommandText:
		if len(parts) < 3 || len(parts) > 6 {
			return "", errors.New("Некорректное количество аргументов")
		}
		var limit int
		isInStock := false
		customerID, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		if len(parts) == 4 {
			if parts[3] != instockArgText {
				return "", errors.New("Некорректные дополнительные аргументы")
			}
			isInStock = true
		}
		if len(parts) == 5 {
			if parts[3] != limitArgText {
				return "", errors.New("Некорректные дополнительные аргументы")
			}
			limit, err = strconv.Atoi(parts[4])
			if err != nil {
				return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
		}
		if len(parts) == 6 {
			limitPos := 3
			if parts[3] == instockArgText {
				limitPos = 4
			} else if parts[5] != instockArgText {
				return "", errors.New("Некорректные дополнительные аргументы")
			}
			isInStock = true
			if parts[limitPos] != limitArgText {
				return "", errors.New("Некорректные дополнительные аргументы")
			}
			limit, err = strconv.Atoi(parts[limitPos+1])
			if err != nil {
				return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
			}
		}

		orders, err := service.GetCustomerOrderList(ctx, pvzID, customerID, limit, isInStock)
		if err != nil {
			return "", fmt.Errorf("Не удалось получить список заказов: %w", err)
		}
		orderTab := ""
		for _, o := range orders {
			orderTab += formOrder(o) + "\n"
		}
		return orderTab, nil
	case takeRefundCommandText:
		if len(parts) != 4 {
			return "", errors.New("Некорректное количество аргументов")
		}
		customerID, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		orderID, err := strconv.Atoi(parts[3])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		err = service.TakeRefundFromCustomer(ctx, pvzID, customerID, orderID)
		if err != nil {
			return "", fmt.Errorf("Не удалось принять возврат от клиента: %w", err)
		}
		return fmt.Sprintf("Возврат заказа %d принят", orderID), nil
	case refundListCommandText:
		if len(parts) != 4 {
			return "", errors.New("Некорректное количество аргументов")
		}
		pageNum, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}
		pageSize, err := strconv.Atoi(parts[3])
		if err != nil {
			return "", fmt.Errorf("Не удалось преобразовать аргумент в число: %w", err)
		}

		orders, err := service.GetRefundList(ctx, pvzID, pageNum, pageSize)
		if err != nil {
			return "", fmt.Errorf("Не удалось получить список возвратов: %w", err)
		}
		orderTab := ""
		for _, o := range orders {
			orderTab += formOrder(o) + "\n"
		}
		return orderTab, nil
	case checkPvzCommandText:
		if len(parts) != 2 {
			return "", errors.New("Некорректное количество аргументов")
		}
		p, err := service.GetPvzByID(ctx, int64(pvzID))
		if err != nil {
			return "", fmt.Errorf("Не удалось получить информацию о ПВЗ: %w", err)
		}
		return fmt.Sprintf("ПВЗ id %d Название: %s Адресс: %s Контакты: %s\n", p.ID, p.Name, p.Adress, p.Contacts), nil

	// тут уникальная логика для консольного ввода, которой не было в задании по http и sql
	// поэтому на месте соберу команду тут, раз она не будет больше нигде использоваться
	// чтобы не добавлять новую функцию в круд ПВЗ
	case newPvzCommandText:
		if len(parts) < 5 {
			return "", errors.New("Некорректное количество аргументов")
		}
		re := regexp.MustCompile(`(".*?")`)
		pvzInput := re.FindAllString(strings.Join(parts[2:], " "), -1)
		if len(pvzInput) != 3 {
			return "", errors.New("Некорректное количество аргументов")
		}
		_, err := service.GetPvzByID(ctx, int64(pvzID))
		if !errors.Is(err, pvz_errors.ErrNotFound) {
			if err != nil {
				return "", fmt.Errorf("Не удалось получить информацию о ПВЗ: %w", err)
			}
			return "", errors.New("ПВЗ уже существует")
		}
		_, err = service.ModifyPvz(ctx, pvz_dto.Pvz{
			ID:       int64(pvzID),
			Name:     strings.Trim(pvzInput[0], `"`),
			Adress:   strings.Trim(pvzInput[1], `"`),
			Contacts: strings.Trim(pvzInput[2], `"`),
		})
		if err != nil {
			return "", fmt.Errorf("Ошибка при записи ПВЗ в базу данных: %w", err)
		}
		return fmt.Sprintf("ПВЗ id %d успешно добавлен в базу.\n", pvzID), nil
	default:
		return "", errors.New("Такой команды нет")
	}
}

func formOrder(o orders_dto.Order) string {
	return fmt.Sprint(
		"id:", o.OrderID, "; customer:", o.CustomerID,
		"; arrived:", o.ArrivalTime.Year(), o.ArrivalTime.Month(), o.ArrivalTime.Day(),
		"; stored until:", o.StorageLastTime.Year(), o.StorageLastTime.Month(), o.StorageLastTime.Day(),
		"; completed:", o.IsCompleted, "; refunded:", o.IsRefunded,
		"; weight:", o.Weight, "; price:", o.Price, "; package:", o.PackageType,
	)
}
