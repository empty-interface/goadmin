package dbms

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GormConnection struct {
	*gorm.DB
	isOpen bool
}
type driver interface {
	open(dsn string) gorm.Dialector
	dsn() string
}

func newGormConnection(driver driver, opts ...gorm.Option) (*GormConnection, error) {
	db, err := gorm.Open(driver.open(driver.dsn()), opts...)
	if err != nil {
		return nil, err
	}
	return &GormConnection{
		DB:     db,
		isOpen: true,
	}, nil
}

func NewClient(driverName string, config config) (*GormConnection, error) {
	driver, ok := supportedDrivers[driverName]
	if !ok {
		return nil, driverNotSupported
	}
	// opt := {}
	return newGormConnection(driver.init(config))
}

type driverInitializer = func(config) driver
type supportedDriver struct {
	Name string
	init driverInitializer
}

var supportedDrivers = map[string]supportedDriver{
	"postgres": {"PostgreSQL", newPostgresDriver},
	"mysql":    {"MySQL", newMySQLDriver},
}

func GetSupportedDrivers() map[string]string {
	drivers := map[string]string{}
	for k, v := range supportedDrivers {
		drivers[k] = v.Name
	}
	return drivers
}

var (
	driverNotSupported = errors.New("Driver Not Supportd")
)

func (conn GormConnection) Close() {
	conn.isOpen = false
}

type config struct {
	Host, Username, DBName, Port, Password string
}

func NewConfig(username, password, dbname, host, port string) config {
	return config{
		Host:     host,
		Username: username,
		DBName:   dbname,
		Port:     port,
		Password: password,
	}
}

type QueryResult struct {
	rows        *sql.Rows
	processTime time.Duration
	err         error
}
