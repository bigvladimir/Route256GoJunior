package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// предполагаю что как работает mux.Vars мне тестировать не надо, поэтому вынес его из тестируемых частей хендлеров
func Test_checkKey(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		key       string
		ok        bool
		wantedInt int64
		wantedErr error
	}{
		{"ok", "1", true, 1, nil},
		{"wrong request", "", false, 0, errors.New("Неправильный ключ запроса")},
		{"key is not int", "one", true, 0, errors.New("Значение ключа не число")},
		{"zero id", "0", true, 0, errors.New("Неположительный id")},
		{"negative id", "-1", true, 0, errors.New("Неположительный id")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualValue, actualErr := checkKey(tc.key, tc.ok)

			assert.Equal(t, tc.wantedInt, actualValue)
			assert.Equal(t, tc.wantedErr, actualErr)
		})
	}
}

func Test_validatePvzReq(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		req       pvzRequest
		wantedErr error
	}{
		{"ok", pvzRequest{"name", "adress", "contacts"}, nil},
		{"missed all arg", pvzRequest{"", "", ""}, errors.New("Пустые поля")},
		{"missed arg 1", pvzRequest{"", "adress", "contacts"}, errors.New("Пустые поля")},
		{"missed arg 2", pvzRequest{"name", "", "contacts"}, errors.New("Пустые поля")},
		{"missed arg 3", pvzRequest{"name", "adress", ""}, errors.New("Пустые поля")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualErr := validatePvzReq(tc.req)

			assert.Equal(t, tc.wantedErr, actualErr)
		})
	}
}

func Test_validateFullPvzReq(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		req       pvzFullRequest
		wantedErr error
	}{
		{"ok", pvzFullRequest{1, "name", "adress", "contacts"}, nil},
		{"zero id", pvzFullRequest{0, "name", "adress", "contacts"}, errors.New("Неположительный id")},
		{"negative id", pvzFullRequest{-1, "name", "adress", "contacts"}, errors.New("Неположительный id")},
		{"missed all str arg", pvzFullRequest{1, "", "", ""}, errors.New("Пустые поля")},
		{"missed str arg 1", pvzFullRequest{1, "", "adress", "contacts"}, errors.New("Пустые поля")},
		{"missed str arg 2", pvzFullRequest{1, "name", "", "contacts"}, errors.New("Пустые поля")},
		{"missed str arg 3", pvzFullRequest{1, "name", "adress", ""}, errors.New("Пустые поля")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualErr := validateFullPvzReq(tc.req)

			assert.Equal(t, tc.wantedErr, actualErr)
		})
	}
}
