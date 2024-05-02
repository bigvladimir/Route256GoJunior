package orders

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"homework/internal/app/orders/dto"
)

//
// блок заказов обновлен просто чтобы работал
// и не мешал основной части 8 дз,
// но вообще, в него нужно добавить транзакции,
// и удалить копипаст валидации
//
// TODO пофиксить возвращение ErrNotFound
//

type storage interface {
	Add(context.Context, dto.OrderInput) error
	Get(context.Context, string, ...int) ([]dto.Order, error)
	Update(context.Context, int, string, string) error
	Delete(context.Context, int) error
}

// Service provides functions for interacting with the order storage
type Service struct {
	stor         storage
	packVariants map[string]PackageVariant
}

// NewService creates Service
func NewService(stor storage, packVariants map[string]PackageVariant) *Service {
	return &Service{
		stor:         stor,
		packVariants: packVariants,
	}
}

// TakeOrderFromCourier обратывает принятие заказа от курьера
func (s *Service) TakeOrderFromCourier(ctx context.Context, order dto.OrderInput) error {
	checkOrder, err := s.get(ctx, "orderID", order.OrderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(checkOrder) != 0 {
		return errors.New("Заказ уже есть в базе")
	}

	if order.PackageType != "" {
		packVariant, ok := s.packVariants[order.PackageType]
		if !ok {
			return errors.New("Не удалось найти упаковку с таким наименованием")
		}
		newOrder, err := packVariant.ApplyPackage(order)
		if err != nil {
			return fmt.Errorf("Ошибка при применении упаковки: %w", err)
		}
		order = newOrder
	}

	order.StorageLastTime = time.Date(
		order.StorageLastTime.Year(), order.StorageLastTime.Month(), order.StorageLastTime.Day(),
		23, 59, 0, 0,
		time.Now().Location(),
	)

	if err := s.add(ctx, order); err != nil {
		return fmt.Errorf("Не удалось добавить заказ в базу данных: %w", err)
	}

	return nil
}

// ReturnOrderToCourier обратывает возврат заказа курьеру,
// если у заказа закончился срок хранения,
// позволяет отдать курьеру заказ возвращенный клиентом независимо от срока хранения
func (s *Service) ReturnOrderToCourier(ctx context.Context, pvzID, orderID int) error {
	orders, err := s.get(ctx, "orderID", orderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(orders) == 0 {
		return errors.New("Заказ не найден")
	}

	order := orders[0]
	if pvzID != order.PvzID {
		return errors.New("Заказ не принадлежит этому ПВЗ")
	}
	if order.IsCompleted && !order.IsRefunded {
		return errors.New("Заказ уже был выдан клиенту")
	}
	if order.StorageLastTime.After(time.Now()) && !order.IsRefunded {
		return errors.New("У заказа ещё не вышел срок хранения")
	}

	err = s.delete(ctx, order.OrderID)
	if err != nil {
		return fmt.Errorf("Не удалось удалить заказ из базы данных: %w", err)
	}

	return nil
}

// GiveOrderToCustomer обрабатывает выдачу заказа
func (s *Service) GiveOrderToCustomer(ctx context.Context, pvzID, customerID, orderID int) error {
	orders, err := s.get(ctx, "orderID", orderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(orders) == 0 {
		return errors.New("Заказ не найден")
	}

	order := orders[0]
	if pvzID != order.PvzID {
		return errors.New("Заказ не принадлежит этому ПВЗ")
	}
	if customerID != order.CustomerID {
		return errors.New("У заказа другой владелец")
	}
	if order.IsCompleted && !order.IsRefunded {
		return errors.New("Попытка забрать уже отданный заказ")
	}
	if order.IsRefunded {
		return errors.New("Попытка забрать возвращенный заказ")
	}
	if time.Now().After(order.StorageLastTime) {
		return errors.New("Недействительный срок хранения")
	}

	if err = s.update(ctx, order.OrderID, "is_completed", "TRUE"); err != nil {
		return fmt.Errorf("Не удалось обновить информацию о заказе в базе данных: %w", err)
	}

	return nil
}

// TakeRefundFromCustomer обрабатывает возврат заказа клиентом
func (s *Service) TakeRefundFromCustomer(ctx context.Context, pvzID, customerID, orderID int) error {
	orders, err := s.get(ctx, "orderID", orderID)
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о заказе: %w", err)
	}
	if len(orders) == 0 {
		return errors.New("Заказ не найден")
	}
	order := orders[0]
	if pvzID != order.PvzID {
		return errors.New("Заказ не принадлежит этому ПВЗ")
	}
	if customerID != order.CustomerID {
		return errors.New("У заказа другой владелец")
	}
	if !order.IsCompleted {
		return errors.New("Заказ найден, но он не выдавался никому")
	}
	if order.IsRefunded {
		return errors.New("Заказ найден, но уже вернули")
	}

	lastTimeForRefund := time.Date(
		order.StorageLastTime.Year(), order.StorageLastTime.Month(), order.StorageLastTime.Day()+2,
		23, 59, 0, 0,
		order.StorageLastTime.Location(),
	)
	if time.Now().After(lastTimeForRefund) {
		return errors.New("Срок возврата истёк")
	}

	err = s.update(ctx, order.OrderID, "is_refunded", "TRUE")
	if err != nil {
		return fmt.Errorf("Не удалось обновить информацию о заказе в базе данных: %w", err)
	}

	return nil
}

// GetRefundList возвращает страницу возвращенных заказов в этом пвз в виде слайса
// pageNum int - номер страницы, pageSize int - размер страницы
func (s *Service) GetRefundList(ctx context.Context, pvzID, pageNum, pageSize int) ([]dto.Order, error) {
	orders, err := s.get(ctx, "isRefunded", pvzID)
	if err != nil {
		err = fmt.Errorf("Не удалось получить информацию о заказах: %w", err)
		return nil, err
	}
	if len(orders) == 0 {
		err := errors.New("Возвраты не найдены")
		return nil, err
	}
	if pageSize*(pageNum-1) >= len(orders) {
		err := errors.New("Возвратов меньше чем запрашиваемая страница")
		return nil, err
	}

	limit := pageSize * (pageNum)
	if cap(orders) < limit {
		limit = cap(orders)
	}
	return orders[pageSize*(pageNum-1) : limit], nil
}

// GetCustomerOrderList возвращает слайс заказов по ID клиента в этом ПВЗ,
// limit int устанавливает максимальное количество возвращаемых заказов,
// если limit = 0, то ограничения нет
// isInStock bool устанавливает необходимость проверки наличия заказ в пункте,
// в том числе возвращенные
func (s *Service) GetCustomerOrderList(ctx context.Context, pvzID, customerID, limit int, isInStock bool) ([]dto.Order, error) {
	if limit < 0 {
		err := errors.New("Отрицательное значение ограничения количества заказов")
		return nil, err
	}
	orders, err := s.get(ctx, "customerID", pvzID, customerID)
	if err != nil {
		err = fmt.Errorf("Не удалось получить информацию о заказах: %w", err)
		return nil, err
	}
	if len(orders) == 0 {
		err := errors.New("Заказы не найдены")
		return nil, err
	}

	var newOrders []dto.Order
	if isInStock {
		for _, order := range orders {
			if !order.IsCompleted || order.IsRefunded {
				newOrders = append(newOrders, order)
			}
		}
		orders = newOrders
	}

	sortOrdersByArrival(orders)
	if limit == 0 || limit > len(orders) {
		limit = len(orders)
	}

	return orders[:limit], nil
}

func (s *Service) add(ctx context.Context, input dto.OrderInput) error {
	if input.OrderID <= 0 {
		return errors.New("Некорректный id заказа")
	}
	if input.CustomerID <= 0 {
		return errors.New("Некорректный id клиента")
	}
	if !input.StorageLastTime.After(time.Now()) {
		return errors.New("Некорректное время хранения заказа")
	}
	if input.Weight <= 0 {
		return errors.New("Некорректный вес заказа")
	}
	if input.Price < 0 {
		return errors.New("Некорректная цена заказа")
	}

	return s.stor.Add(ctx, input)
}

func (s *Service) get(ctx context.Context, filter string, id ...int) ([]dto.Order, error) {
	if len(filter) == 0 {
		err := errors.New("Пустой фильтр для Get()")
		return nil, err
	}
	for _, j := range id {
		if j <= 0 {
			err := errors.New("Некорректный id для Get()")
			return nil, err
		}
	}

	return s.stor.Get(ctx, filter, id...)
}

func (s *Service) update(ctx context.Context, id int, column, value string) error {
	if id <= 0 {
		return errors.New("Некорректный id заказа при обновлении")
	}

	return s.stor.Update(ctx, id, column, value)
}

func (s *Service) delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("Некорректный id заказа")
	}

	return s.stor.Delete(ctx, id)
}

func sortOrdersByArrival(orders []dto.Order) {
	sort.Slice(orders, func(i, j int) bool {
		return orders[j].ArrivalTime.After(orders[i].ArrivalTime)
	})
}
