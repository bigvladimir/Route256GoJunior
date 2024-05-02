package server

import (
	"context"
	"errors"
	"homework/internal/app/pvz"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)

	tt := []struct {
		name         string
		getReturnPvz pvz.Pvz
		getReturnErr error
		wantedJSON   string
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", pvz.Pvz{1, "a", "b", "c"}, nil, `{"id":1,"name":"a","adress":"b","contacts":"c"}`, http.StatusOK, nil},
		{"not found", pvz.Pvz{}, pvz.ErrNotFound, "", http.StatusNotFound, pvz.ErrNotFound},
		{"internal error", pvz.Pvz{}, errors.New("Any internal err"), "", http.StatusInternalServerError, errors.New("Any internal err")},
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

func Test_Delete(t *testing.T) {
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
		{"not found", pvz.ErrNotFound, http.StatusNotFound, pvz.ErrNotFound},
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

func Test_Create(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.Background()
		req   = pvzRequest{"name", "adress", "contacts"}
		bdReq = pvz.PvzInput{"name", "adress", "contacts"}
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

func Test_Modify(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.Background()
		req   = pvzFullRequest{1, "name", "adress", "contacts"}
		bdReq = pvz.Pvz{1, "name", "adress", "contacts"}
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

func Test_Update(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.Background()
		id    = int64(1)
		req   = pvzRequest{"name", "adress", "contacts"}
		bdReq = pvz.Pvz{1, "name", "adress", "contacts"}
	)

	tt := []struct {
		name         string
		getReturnErr error
		wantedStatus int
		wantedErr    error
	}{
		{"smoke", nil, http.StatusOK, nil},
		{"not found", pvz.ErrNotFound, http.StatusNotFound, pvz.ErrNotFound},
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
