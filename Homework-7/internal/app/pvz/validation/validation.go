package validation

import (
	"errors"

	"homework/internal/app/pvz/dto"
)

// ValidatePvz checks Pvz ID is correct and fields are not empty
func ValidatePvz(pvzObj dto.Pvz) error {
	if pvzObj.ID <= 0 {
		return errors.New("Incorrect ID")
	}
	if pvzObj.Name == "" || pvzObj.Adress == "" || pvzObj.Contacts == "" {
		return errors.New("Empty field")
	}

	return nil
}

// ValidatePvzInput checks PvzInput fields are not empty
func ValidatePvzInput(pvzObj dto.PvzInput) error {
	if pvzObj.Name == "" || pvzObj.Adress == "" || pvzObj.Contacts == "" {
		return errors.New("Empty field")
	}

	return nil
}

// ValidatePvzID checks Pvz ID is correct
func ValidatePvzID(id int64) error {
	if id <= 0 {
		return errors.New("Incorrect ID")
	}

	return nil
}
