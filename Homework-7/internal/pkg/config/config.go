package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var errUninitializedConfig = errors.New("Config not initialized")

var cfg config

type bdParams struct {
	Host     string `yaml:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT"`
	DBname   string `yaml:"POSTGRES_DB"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
}

type config struct {
	initialized bool

	bd          bdParams
	brokerPorts []string
	users       map[string]string
}

// Init is a necessary initialization function if you want to use config
func Init() error {
	if err := cfg.setBd(); err != nil {
		return fmt.Errorf("setBd: %w", err)
	}
	if err := cfg.setBrokerPorts(); err != nil {
		return fmt.Errorf("GetBrokerPorts: %w", err)
	}
	if err := cfg.setUsers(); err != nil {
		return fmt.Errorf("GetUsers: %w", err)
	}

	cfg.initialized = true

	return nil
}

// Cfg checks initialization and return config
func Cfg() *config {
	if !cfg.initialized {
		panic(errUninitializedConfig)
	}
	return &cfg
}

// GetBdDSN returns bd login string
func (c *config) GetBdDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.bd.Host, c.bd.Port, c.bd.User, c.bd.Password, c.bd.DBname,
	)
}

// GetBdDSN returns kafka brokers hosts
func (c *config) GetBrokerPorts() []string {
	brokerPortsCopy := make([]string, len(c.brokerPorts))
	copy(brokerPortsCopy, c.brokerPorts)
	return brokerPortsCopy
}

// GetUserPassword returns specific user password,
// not very safe =)
func (c *config) GetUserPassword(username string) (string, bool) {
	v, ok := c.users[username]
	return v, ok
}

func (c *config) setBd() error {
	readedData, err := os.ReadFile("./config/bdcfg.yaml")
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	bd := bdParams{}
	err = yaml.Unmarshal(readedData, &bd)
	if err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	c.bd = bd

	return nil
}

func (c *config) setBrokerPorts() error {
	yamlFile, err := os.ReadFile("./config/brokerports.yaml")
	if err != nil {
		return fmt.Errorf("io.ReadFile: %w", err)
	}
	data := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	var brokers []string
	for _, value := range data {
		brokers = append(brokers, value)
	}

	c.brokerPorts = brokers

	return nil
}

func (c *config) setUsers() error {
	yamlFile, err := os.ReadFile("./config/users.yaml")
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	users := make(map[string]string)
	if err := yaml.Unmarshal(yamlFile, &users); err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	c.users = users

	return nil
}
