package grpcserver

import (
	"errors"
	"time"
)

//
// возможно правильнее было бы идентифицировать ошибку валидации приходящую изнутри кода,
// но там всё так неорганизованно, что проще сразу провалидировать,
// чтобы вернуть правильный код ответа
//

// orders structs begin

func (m *orderIdentifierModel) validate() error {
	if m.orderID <= 0 {
		return errors.New("Incorrect order ID")
	}
	if m.pvzID <= 0 {
		return errors.New("Incorrect pvz ID")
	}
	return nil
}

func (m *orderBodyModel) validate() error {
	if m.customerID <= 0 {
		return errors.New("Incorrect customer ID")
	}
	if m.weight <= 0 {
		return errors.New("Incorrect weight")
	}
	if m.price < 0 {
		return errors.New("Incorrect price")
	}
	if !m.storageLastTime.After(time.Now()) {
		return errors.New("Incorrect storage time")
	}

	return nil
}

func (m *orderInfoModel) validate() error {
	if err := m.identifier.validate(); err != nil {
		return nil
	}
	if err := m.body.validate(); err != nil {
		return nil
	}

	return nil
}

func (m *orderRequestForCustomerModel) validate() error {
	if err := m.identifier.validate(); err != nil {
		return nil
	}
	if m.customerID <= 0 {
		return errors.New("Incorrect customer ID")
	}

	return nil
}

func (m *refundListRequestModel) validate() error {
	if m.pvzID <= 0 {
		return errors.New("Incorrect pvz ID")
	}
	if m.pageNum < 1 {
		return errors.New("Incorrect page number")
	}
	if m.pageSize < 1 {
		return errors.New("Incorrect page size")
	}

	return nil
}

func (m *customerOrderListRequestModel) validate() error {
	if m.pvzID <= 0 {
		return errors.New("Incorrect pvz ID")
	}
	if m.customerID <= 0 {
		return errors.New("Incorrect customer ID")
	}
	if m.limit < 0 {
		return errors.New("Incorrect limit")
	}

	return nil
}

// orders structs end

// pvz structs begin

func (m *pvzIdentifierModel) validate() error {
	if m.pvzID <= 0 {
		return errors.New("Incorrect pvz ID")
	}

	return nil
}

func (m *pvzInfoModel) validate() error {
	if m.name == "" || m.adress == "" || m.contacts == "" {
		return errors.New("Contains empty fields")
	}

	return nil
}

func (m *pvzModel) validate() error {
	if err := m.identifier.validate(); err != nil {
		return nil
	}
	if err := m.info.validate(); err != nil {
		return nil
	}

	return nil
}

// pvz structs end
