package orders

import (
	"errors"

	"homework/internal/app/orders/dto"
)

const (
	// BagVariantName string name for NewBagPackage
	BagVariantName string = "bag"
	// BoxVariantName string name for NewBoxPackage
	BoxVariantName string = "box"
	// FilmVariantName string name for NewFilmPackage
	FilmVariantName string = "film"
)

// PackageVariant contains the order packaging characteristics
type PackageVariant struct {
	maxWeight float64
	price     int
}

// ApplyPackage checks PackageVariant restrictions ans applies properties
func (pv PackageVariant) ApplyPackage(order dto.OrderInput) (dto.OrderInput, error) {
	if order.Weight <= 0 {
		return dto.OrderInput{}, errors.New("Неположительный вес заказа")
	}
	if order.Price < 0 {
		return dto.OrderInput{}, errors.New("Отрицательная цена заказа")
	}

	if pv.maxWeight > 0 && order.Weight > pv.maxWeight {
		return dto.OrderInput{}, errors.New("Заказ слишком тяжелый")
	}
	order.Price += pv.price
	return order, nil
}

// NewBagPackage creates PackageVariant with properties: maxWeight = 10.0, price = 5
func NewBagPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 10.0,
		price:     5,
	}
}

// NewBoxPackage creates PackageVariant with properties: maxWeight = 30.0, price = 20
func NewBoxPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 30.0,
		price:     20,
	}
}

// NewFilmPackage creates PackageVariant with properties: maxWeight = unlimited(0), price = 1
func NewFilmPackage() PackageVariant {
	return PackageVariant{
		maxWeight: 0.0,
		price:     1,
	}
}
