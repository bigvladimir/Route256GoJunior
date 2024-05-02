package postrgesql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"homework/internal/app/pvz/dto"
	pvz_errors "homework/internal/app/pvz/errors"
	"homework/internal/app/pvz/repository"
	"homework/internal/pkg/db/transaction_manager"
)

// не совсем понимаю как можно убрать зависимость от этого интерфейса
type queryEngineManager interface {
	GetQueryEngine(ctx context.Context) transaction_manager.DbOps
}

// PvzStore provides functions for the direct interacting with the pvz database
type PvzStore struct {
	qe queryEngineManager
}

// NewPvzStorage creates PvzStore
func NewPvzStorage(queryEngine queryEngineManager) *PvzStore {
	return &PvzStore{qe: queryEngine}
}

// Add добавляет запись без указания id
func (r *PvzStore) Add(ctx context.Context, input dto.PvzInput) (int64, error) {
	var pvz repository.PvzInputRow
	pvz.MapFromDTO(input)
	var id int64
	err := r.qe.GetQueryEngine(ctx).ExecQueryRow(
		ctx, "INSERT INTO pvz(name, adress, contacts) VALUES ($1, $2, $3) RETURNING id;", pvz.Name, pvz.Adress, pvz.Contacts,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("ExecQueryRow INSERT: %w", err)
		return 0, err
	}
	return id, nil
}

// GetByID возвращает строку по id, если она существует
func (r *PvzStore) GetByID(ctx context.Context, id int64) (dto.Pvz, error) {
	var pvz repository.PvzRow
	err := r.qe.GetQueryEngine(ctx).Get(ctx, &pvz, "SELECT id, name, adress, contacts FROM pvz where id=$1;", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Pvz{}, pvz_errors.ErrNotFound
		}
		err = fmt.Errorf("Get SELECT: %w", err)
		return dto.Pvz{}, err
	}
	return pvz.MapToDTO(), nil
}

// Update обновляет строку по id, если она существует
func (r *PvzStore) Update(ctx context.Context, input dto.Pvz) error {
	var pvz repository.PvzRow
	pvz.MapFromDTO(input)
	commandTag, err := r.qe.GetQueryEngine(ctx).Exec(
		ctx, "UPDATE pvz SET name = $1, adress = $2, contacts = $3 WHERE id = $4;",
		pvz.Name, pvz.Adress, pvz.Contacts, pvz.ID,
	)
	if err != nil {
		err = fmt.Errorf("Exec UPDATE: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return pvz_errors.ErrNotFound
	}
	return nil
}

// Delete удаляет строку по id, если она существует
func (r *PvzStore) Delete(ctx context.Context, id int64) error {
	commandTag, err := r.qe.GetQueryEngine(ctx).Exec(ctx, "DELETE FROM pvz WHERE id = $1;", id)
	if err != nil {
		err = fmt.Errorf("Exec DELETE: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return pvz_errors.ErrNotFound
	}
	return nil
}

// Modify обновляет данные по id, если не находит, что обновить, то вставляет новые
func (r *PvzStore) Modify(ctx context.Context, pvz dto.Pvz) (int64, error) {
	err := r.Update(ctx, pvz)
	if errors.Is(err, pvz_errors.ErrNotFound) {
		var id int64
		err = r.qe.GetQueryEngine(ctx).ExecQueryRow(
			ctx, "INSERT INTO pvz(id, name, adress, contacts) VALUES ($1, $2, $3, $4) RETURNING id;", pvz.ID, pvz.Name, pvz.Adress, pvz.Contacts,
		).Scan(&id)
		if err != nil {
			err = fmt.Errorf("ExecQueryRow INSERT: %w", err)
			return 0, err
		}
		return id, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, nil
}
