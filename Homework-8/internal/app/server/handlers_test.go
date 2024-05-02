package server

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pvz_dto "homework/internal/app/pvz/dto"
	pvz_errors "homework/internal/app/pvz/errors"
)

func TestGetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)

	tt := []struct {
		name         string
		getReturnPvz pvz_dto.Pvz
		getReturnErr error
		wantedJSON   string
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", pvz_dto.Pvz{ID: 1, Name: "a", Adress: "b", Contacts: "c"}, nil, `{"id":1,"name":"a","adress":"b","contacts":"c"}`, http.StatusOK, nil},
		{"not found", pvz_dto.Pvz{}, pvz_errors.ErrNotFound, "", http.StatusNotFound, pvz_errors.ErrNotFound},
		{"internal error", pvz_dto.Pvz{}, errors.New("Any internal err"), "", http.StatusInternalServerError, errors.New("Any internal err")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCoreOps.EXPECT().GetPvzByID(gomock.Any(), id).Return(tc.getReturnPvz, tc.getReturnErr)

			result, status, err := s.srv.getByIDexec(ctx, id)

			assert.Equal(t, tc.wantedJSON, string(result))
			assert.Equal(t, tc.wantedStatus, status)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)

	tt := []struct {
		name         string
		getReturnErr error
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", nil, http.StatusOK, nil},
		{"not found", pvz_errors.ErrNotFound, http.StatusNotFound, pvz_errors.ErrNotFound},
		{"internal error", errors.New("Any internal err"), http.StatusInternalServerError, errors.New("Any internal err")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCoreOps.EXPECT().DeletePvz(gomock.Any(), id).Return(tc.getReturnErr)

			status, err := s.srv.deleteExec(ctx, id)

			assert.Equal(t, tc.wantedStatus, status)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		req = pvzRequest{
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
		bdReq = pvz_dto.PvzInput{
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
	)

	tt := []struct {
		name         string
		getReturnID  int64
		getReturnErr error
		wantedJSON   string
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", 1, nil, `{"id":1,"name":"name","adress":"adress","contacts":"contacts"}`, http.StatusCreated, nil},
		{"internal error", 0, errors.New("Any internal err"), "", http.StatusInternalServerError, errors.New("Any internal err")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCoreOps.EXPECT().AddPvz(gomock.Any(), bdReq).Return(tc.getReturnID, tc.getReturnErr)

			result, status, err := s.srv.createExec(ctx, req)

			assert.Equal(t, tc.wantedJSON, string(result))
			assert.Equal(t, tc.wantedStatus, status)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}

func TestModify(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		req = pvzFullRequest{
			ID:       1,
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
		bdReq = pvz_dto.Pvz{
			ID:       1,
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
	)

	tt := []struct {
		name         string
		getReturnID  int64
		getReturnErr error
		wantedJSON   string
		wantedStatus int
		wantedErr    error
	}{
		{"updated", 0, nil, `{"id":1,"name":"name","adress":"adress","contacts":"contacts"}`, http.StatusOK, nil},
		{"created", 1, nil, `{"id":1,"name":"name","adress":"adress","contacts":"contacts"}`, http.StatusCreated, nil},
		{"internal error", 0, errors.New("Any internal err"), "", http.StatusInternalServerError, errors.New("Any internal err")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCoreOps.EXPECT().ModifyPvz(gomock.Any(), bdReq).Return(tc.getReturnID, tc.getReturnErr)

			result, status, err := s.srv.modifyExec(ctx, req)

			assert.Equal(t, tc.wantedJSON, string(result))
			assert.Equal(t, tc.wantedStatus, status)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
		req = pvzRequest{
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
		bdReq = pvz_dto.Pvz{
			ID:       1,
			Name:     "name",
			Adress:   "adress",
			Contacts: "contacts",
		}
	)

	tt := []struct {
		name         string
		getReturnErr error
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", nil, http.StatusOK, nil},
		{"not found", pvz_errors.ErrNotFound, http.StatusNotFound, pvz_errors.ErrNotFound},
		{"internal error", errors.New("Any internal err"), http.StatusInternalServerError, errors.New("Any internal err")},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCoreOps.EXPECT().UpdatePvz(gomock.Any(), bdReq).Return(tc.getReturnErr)

			status, err := s.srv.updateExec(ctx, id, req)

			assert.Equal(t, tc.wantedStatus, status)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
