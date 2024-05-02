package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v2"

	"homework/internal/pkg/utils"
)

// NewDb creates new database connection
func NewDb(ctx context.Context, dockerComposeFile string) (*Database, error) {
	dsn, err := generateDsn(dockerComposeFile)
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

type dockerComposeBd struct {
	Services struct {
		Postgres struct {
			Image string `yaml:"image"`
			Env   struct {
				Db       string `yaml:"POSTGRES_DB"`
				User     string `yaml:"POSTGRES_USER"`
				Password string `yaml:"POSTGRES_PASSWORD"`
			} `yaml:"environment"`
			Ports []string `yaml:"ports"`
		} `yaml:"postgres"`
	} `yaml:"services"`
}

func generateDsn(dockerComposeFile string) (string, error) {
	rootPath, err := utils.GetRootPath()
	if err != nil {
		return "", fmt.Errorf("utils.GetRootPath: %w", err)
	}
	dockerComposeFile = filepath.Join(rootPath, filepath.Clean(dockerComposeFile))
	readedData, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return "", fmt.Errorf("Ошибка при чтении конфига бд: %w", err)
	}
	dockerCompose := dockerComposeBd{}
	err = yaml.Unmarshal(readedData, &dockerCompose)
	if err != nil {
		return "", fmt.Errorf("Ошибка при записи конфига в структуру: %w", err)
	}

	db := dockerCompose.Services.Postgres.Env.Db
	user := dockerCompose.Services.Postgres.Env.User
	password := dockerCompose.Services.Postgres.Env.Password
	ports := dockerCompose.Services.Postgres.Ports
	port := strings.Split(ports[0], ":")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		hostCfg, port[0], user, password, db)

	return dsn, nil
}
