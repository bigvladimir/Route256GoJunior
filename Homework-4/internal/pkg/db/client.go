package db

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/pkg/jsonutil"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context) (*Database, error) {
	dsn, err := generateDsn()
	if err != nil {
		err = fmt.Errorf("Ошибка при генерации запроса подключения к базе данных: %w", err)
		return nil, err
	}
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		err = fmt.Errorf("Ошибка при подключении к базе данных: %w", err)
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn() (string, error) {

	readedData, err := jsonutil.ReadJSON(make(map[string]interface{}), bdCfgPAth)
	if err != nil {
		err = fmt.Errorf("Ошибка при чтении конфига бд: %w", err)
		return "", err
	}
	castData := readedData.(map[string]interface{})

	dbCfgNames := []string{
		hostCfg,
		portCfg,
		userCfg,
		passwordCfg,
		dbnameCfg,
	}
	data := make(map[string]string)
	for _, name := range dbCfgNames {
		s, ok := castData[name]
		if !ok {
			err := errors.New(cfgNotFoundErr + name)
			return "", err
		}
		data[name] = fmt.Sprint(s)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		data[hostCfg], data[portCfg], data[userCfg], data[passwordCfg], data[dbnameCfg])

	return dsn, nil
}
