//go:build integration

package tests

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pvz_dto "homework/internal/app/pvz/dto"
)

const pvzURL = "http://anyhost:1234/pvz"

func fillPvzDb(t *testing.T, pvz pvz_dto.Pvz) {
	_, err := db.DB.Exec(
		context.Background(),
		`INSERT INTO pvz(id,name,adress,contacts) VALUES ($1,$2,$3,$4);`,
		pvz.ID, pvz.Name, pvz.Adress, pvz.Contacts,
	)
	require.NoError(t, err)
}

func TestCreate(t *testing.T) {
	tt := []struct {
		name         string
		reqBody      string
		wantedStatus int
	}{
		{
			name:         "created",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusCreated,
		},
		{
			name:         "empty data 1",
			reqBody:      `{"name":"","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 2",
			reqBody:      `{"name":"a","adress":"","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 3",
			reqBody:      `{"name":"a","adress":"b","contacts":""}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "not full json",
			reqBody:      `{"name":"a","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty json",
			reqBody:      "",
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "bad json",
			reqBody:      `{name:a, adress:b, contacts:c}`,
			wantedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.SetUp(t, "pvz")
			defer db.TearDown(t)

			req, err := http.NewRequest(http.MethodPost, pvzURL, strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			defer req.Body.Close()
			rr := httptest.NewRecorder()

			serv.H.Create(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantedStatus, resp.StatusCode)
		})
	}
}

func TestModify(t *testing.T) {
	tt := []struct {
		name         string
		reqBody      string
		wantedStatus int
	}{
		{
			name:         "created",
			reqBody:      `{"id":1,"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusCreated,
		},
		{
			name:         "zero id",
			reqBody:      `{"id":0,"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "negative id",
			reqBody:      `{"id":-1,"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 1",
			reqBody:      `{"id":1,"name":"","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 2",
			reqBody:      `{"id":1,"name":"a","adress":"","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 3",
			reqBody:      `{"id":1,"name":"a","adress":"b","contacts":""}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "not full json",
			reqBody:      `{"id":1,"adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty json",
			reqBody:      "",
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "bad json",
			reqBody:      `{id:1, name:a, adress:b, contacts:c}`,
			wantedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.SetUp(t, "pvz")
			defer db.TearDown(t)

			req, err := http.NewRequest(http.MethodPut, pvzURL, strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			defer req.Body.Close()
			rr := httptest.NewRecorder()

			serv.H.Modify(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantedStatus, resp.StatusCode)
		})
	}
	t.Run("ok", func(t *testing.T) {
		db.SetUp(t, "pvz")
		defer db.TearDown(t)
		fillPvzDb(t, pvz_dto.Pvz{ID: 1, Name: "any", Adress: "any", Contacts: "any"})

		reqBody := `{"id":1,"name":"a","adress":"b","contacts":"c"}`
		wantedStatus := http.StatusOK

		req, err := http.NewRequest(http.MethodPut, pvzURL, strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()
		rr := httptest.NewRecorder()

		serv.H.Modify(rr, req)
		resp := rr.Result()
		defer resp.Body.Close()

		assert.Equal(t, wantedStatus, resp.StatusCode)
	})
}

func TestGetByID(t *testing.T) {
	tt := []struct {
		name         string
		postfixURL   string
		wantedStatus int
	}{
		{name: "ok", postfixURL: "1", wantedStatus: http.StatusOK},
		{name: "not found", postfixURL: "2", wantedStatus: http.StatusNotFound},
		{name: "zero id", postfixURL: "0", wantedStatus: http.StatusBadRequest},
		{name: "negative id", postfixURL: "-1", wantedStatus: http.StatusBadRequest},
		{name: "empty id", postfixURL: "", wantedStatus: http.StatusBadRequest},
		{name: "not int id", postfixURL: "one", wantedStatus: http.StatusBadRequest},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.SetUp(t, "pvz")
			defer db.TearDown(t)
			fillPvzDb(t, pvz_dto.Pvz{ID: 1, Name: "any", Adress: "any", Contacts: "any"})

			req, err := http.NewRequest(http.MethodGet, pvzURL+"/"+tc.postfixURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": tc.postfixURL})
			rr := httptest.NewRecorder()

			serv.H.GetByID(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantedStatus, resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				respBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				assert.Equal(t, `{"id":1,"name":"any","adress":"any","contacts":"any"}`, string(respBody))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tt := []struct {
		name         string
		postfixURL   string
		wantedStatus int
	}{
		{name: "ok", postfixURL: "1", wantedStatus: http.StatusOK},
		{name: "not found", postfixURL: "2", wantedStatus: http.StatusNotFound},
		{name: "zero id", postfixURL: "0", wantedStatus: http.StatusBadRequest},
		{name: "negative id", postfixURL: "-1", wantedStatus: http.StatusBadRequest},
		{name: "empty id", postfixURL: "", wantedStatus: http.StatusBadRequest},
		{name: "not int id", postfixURL: "one", wantedStatus: http.StatusBadRequest},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.SetUp(t, "pvz")
			defer db.TearDown(t)
			fillPvzDb(t, pvz_dto.Pvz{ID: 1, Name: "any", Adress: "any", Contacts: "any"})

			req, err := http.NewRequest(http.MethodDelete, pvzURL+"/"+tc.postfixURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": tc.postfixURL})
			rr := httptest.NewRecorder()

			serv.H.Delete(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantedStatus, resp.StatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	tt := []struct {
		name         string
		postfixURL   string
		reqBody      string
		wantedStatus int
	}{
		{
			name:         "ok",
			postfixURL:   "1",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusOK,
		},
		{
			name:         "not found",
			postfixURL:   "2",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "zero id",
			postfixURL:   "0",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "negative id",
			postfixURL:   "-1",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty id",
			postfixURL:   "",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "not int id",
			postfixURL:   "one",
			reqBody:      `{"name":"a","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 1",
			postfixURL:   "1",
			reqBody:      `{"name":"","adress":"b","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 2",
			postfixURL:   "1",
			reqBody:      `{"name":"a","adress":"","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty data 3",
			postfixURL:   "1",
			reqBody:      `{"name":"a","adress":"b","contacts":""}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "not full json",
			postfixURL:   "1",
			reqBody:      `{"name":"a","contacts":"c"}`,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "empty json",
			postfixURL:   "1",
			reqBody:      "",
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "bad json",
			postfixURL:   "1",
			reqBody:      `{name:a adress:b contacts:c}`,
			wantedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db.SetUp(t, "pvz")
			defer db.TearDown(t)
			fillPvzDb(t, pvz_dto.Pvz{ID: 1, Name: "any", Adress: "any", Contacts: "any"})

			req, err := http.NewRequest(http.MethodPatch, pvzURL+"/"+tc.postfixURL, strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": tc.postfixURL})
			rr := httptest.NewRecorder()

			serv.H.Update(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.wantedStatus, resp.StatusCode)
		})
	}
}
