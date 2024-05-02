package orders

import (
	"errors"
)

const (
	BagVariantName  string = "bag"
	BoxVariantName  string = "box"
	FilmVariantName string = "film"
)

type PackageVariant struct {
	maxWeight float64
	price     int
}

func (pv PackageVariant) ApplyPackage(order OrderInput) (OrderInput, error) {
	if order.Weight <= 0 {
		return OrderInput{}, errors.New("Неположительный вес заказа")
	}
	if order.Price < 0 {
		return OrderInput{}, errors.New("Отрицательная цена заказа")
	}

	if pv.maxWeight > 0 && order.Weight > pv.maxWeight {
		return OrderInput{}, errors.New("Заказ слишком тяжелый")
	}
	order.Price += pv.price
	return order, nil
}

func NewBagPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 10.0,
		price:     5,
	}
}

func NewBoxPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 30.0,
		price:     20,
	}
}

func NewFilmPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 0.0,
		price:     1,
	}
}
