package postrgesql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"homework/internal/app/pvz/dto"
	pvz_errors "homework/internal/app/pvz/errors"
	"homework/internal/app/pvz/repository"
)

type dbOps interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetPool(_ context.Context) *pgxpool.Pool
}

type PvzStore struct {
	db dbOps
}

// NewPvzStorage инициализация базы данных
func NewPvzStorage(database dbOps) *PvzStore {
	return &PvzStore{db: database}
}

func (r *PvzStore) Add(ctx context.Context, input dto.PvzInput) (int64, error) {
	var pvz repository.PvzInputRow
	pvz.MapFromDTO(input)
	var id int64
	err := r.db.ExecQueryRow(
		ctx, "INSERT INTO pvz(name, adress, contacts) VALUES ($1, $2, $3) RETURNING id;", pvz.Name, pvz.Adress, pvz.Contacts,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("Ошибка при записи данных: %w", err)
		return 0, err
	}
	return id, nil
}

func (r *PvzStore) GetByID(ctx context.Context, id int64) (dto.Pvz, error) {
	var pvz repository.PvzRow
	err := r.db.Get(ctx, &pvz, "SELECT id, name, adress, contacts FROM pvz where id=$1;", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Pvz{}, pvz_errors.ErrNotFound
		}
		err = fmt.Errorf("Ошибка при получении данных: %w", err)
		return dto.Pvz{}, err
	}
	return pvz.MapToDTO(), nil
}

func (r *PvzStore) Update(ctx context.Context, input dto.Pvz) error {
	var pvz repository.PvzRow
	pvz.MapFromDTO(input)
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
		return pvz_errors.ErrNotFound
	}
	return nil
}

func (r *PvzStore) Delete(ctx context.Context, id int64) error {
	commandTag, err := r.db.Exec(ctx, "DELETE FROM pvz WHERE id = $1;", id)
	if err != nil {
		err = fmt.Errorf("Ошибка при удалении данных: %w", err)
		return err
	}
	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return pvz_errors.ErrNotFound
	}
	return nil
}

func (r *PvzStore) Modify(ctx context.Context, pvz dto.Pvz) (int64, error) {
	err := r.Update(ctx, pvz)
	if errors.Is(err, pvz_errors.ErrNotFound) {
		var id int64
		err = r.db.ExecQueryRow(
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
