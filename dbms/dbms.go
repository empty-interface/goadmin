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
	// a, err := gorm.Open()
	// a.Delete()
	return &GormConnection{
		DB:     db,
		isOpen: true,
	}, nil
}

func New(driver string, config Config) (*GormConnection, error) {
	_func, ok := supportedDrivers[driver]
	if !ok {
		return nil, driverNotSupported
	}
	// opt := {}
	return newGormConnection(_func(config))
}

type driverInitializer = func(Config) driver

var supportedDrivers = map[string]driverInitializer{
	"postgres": newPostgresDriver,
}

var (
	driverNotSupported = errors.New("Driver Not Supportd")
)

func (conn GormConnection) Close() {
	conn.isOpen = false
}

type Config struct {
	Host, Username, DBName, Port, Password string
}
type QueryResult struct {
	rows        *sql.Rows
	processTime time.Duration
	err         error
}
