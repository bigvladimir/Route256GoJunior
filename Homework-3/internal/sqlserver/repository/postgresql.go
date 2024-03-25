package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"homework/internal/sqlserver/db"
)

type PvzStore struct {
	db *db.Database
}

// NewPvzStore инициализация базы данных
func NewPvzStore(database *db.Database) *PvzStore {
	return &PvzStore{db: database}
}

// Add добавляет запись без указания id
func (r *PvzStore) Add(ctx context.Context, pvz *PvzDTO) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx, "INSERT INTO pvz(name, adress, contacts) VALUES ($1, $2, $3) RETURNING id;", pvz.Name, pvz.Adress, pvz.Contacts,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("Ошибка при записи данных: %w", err)
		return -1, err
	}
	return id, nil
}

// GetByID возвращает строку по id, если она существует
func (r *PvzStore) GetByID(ctx context.Context, id int64) (*PvzDTO, error) {
	var pvz PvzDTO
	err := r.db.Get(ctx, &pvz, "SELECT id, name, adress, contacts FROM pvz where id=$1;", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		err = fmt.Errorf("Ошибка при получении данных: %w", err)
		return nil, err
	}
	return &pvz, nil
}

// Update обновляет строку по id, если она существует,
// TODO добавить частичное обновление
func (r *PvzStore) Update(ctx context.Context, pvz *PvzDTO) error {
	commandTag, err := r.db.Exec(
		ctx, "UPDATE pvz SET name = $1, adress = $2, contacts = $3 WHERE id = $4;",
		pvz.Name, pvz.Adress, pvz.Contacts, pvz.ID,
	)
	if err != nil {
		err = fmt.Errorf("Ошибка при обновлении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete удаляет строку по id, если она существует
func (r *PvzStore) Delete(ctx context.Context, id int64) error {
	commandTag, err := r.db.Exec(ctx, "DELETE FROM pvz WHERE id = $1;", id)
	if err != nil {
		err = fmt.Errorf("Ошибка при удалении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Modify обновляет данные по id, если не находит, что обновить, то вставляет новые
func (r *PvzStore) Modify(ctx context.Context, pvz *PvzDTO) (int64, error) {
	err := r.Update(ctx, pvz)
	if errors.Is(err, ErrNotFound) {
		var id int64
		err := r.db.ExecQueryRow(
			ctx, "INSERT INTO pvz(id, name, adress, contacts) VALUES ($1, $2, $3, $4) RETURNING id;", pvz.ID, pvz.Name, pvz.Adress, pvz.Contacts,
		).Scan(&id)
		if err != nil {
			err = fmt.Errorf("Ошибка при записи данных: %w", err)
			return -1, err
		}
		return id, nil
	}
	if err != nil {
		return -1, err
	}
	return -1, nil
}
