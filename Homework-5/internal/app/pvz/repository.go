package pvz

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/pkg/db"

	"github.com/jackc/pgx/v4"
)

type PvzStore struct {
	db db.DBops
}

type storage interface {
	add(context.Context, PvzInput) (int64, error)
	getByID(context.Context, int64) (Pvz, error)
	update(context.Context, Pvz) error
	delete(context.Context, int64) error
	modify(context.Context, Pvz) (int64, error)
}

// NewPvzStore инициализация базы данных
func NewPvzStorage(database db.DBops) *PvzStore {
	return &PvzStore{db: database}
}

// Add добавляет запись без указания id
func (r *PvzStore) add(ctx context.Context, input PvzInput) (int64, error) {
	pvz := input.mapToModel()
	var id int64
	err := r.db.ExecQueryRow(
		ctx, "INSERT INTO pvz(name, adress, contacts) VALUES ($1, $2, $3) RETURNING id;", pvz.name, pvz.adress, pvz.contacts,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("Ошибка при записи данных: %w", err)
		return 0, err
	}
	return id, nil
}

// GetByID возвращает строку по id, если она существует
func (r *PvzStore) getByID(ctx context.Context, id int64) (Pvz, error) {
	var pvz pvzRow
	err := r.db.Get(ctx, &pvz, "SELECT id, name, adress, contacts FROM pvz where id=$1;", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Pvz{}, ErrNotFound
		}
		err = fmt.Errorf("Ошибка при получении данных: %w", err)
		return Pvz{}, err
	}
	return pvz.mapToDTO(), nil
}

// Update обновляет строку по id, если она существует,
// TODO добавить частичное обновление
func (r *PvzStore) update(ctx context.Context, input Pvz) error {
	pvz := input.mapToModel()
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
func (r *PvzStore) delete(ctx context.Context, id int64) error {
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
func (r *PvzStore) modify(ctx context.Context, pvz Pvz) (int64, error) {
	err := r.update(ctx, pvz)
	if errors.Is(err, ErrNotFound) {
		var id int64
		err := r.db.ExecQueryRow(
			ctx, "INSERT INTO pvz(id, name, adress, contacts) VALUES ($1, $2, $3, $4) RETURNING id;", pvz.ID, pvz.Name, pvz.Adress, pvz.Contacts,
		).Scan(&id)
		if err != nil {
			err = fmt.Errorf("Ошибка при записи данных: %w", err)
			return 0, err
		}
		return id, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, nil
}
