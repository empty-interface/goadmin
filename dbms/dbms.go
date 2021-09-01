package dbms

import (
	"database/sql"
	"errors"
	"fmt"
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
type Table struct {
}

func (conn *GormConnection) GetTables() {
	query := `select * from information_schema.tables where table_schema = public`
	// tables := make([]Table, 0)
	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return
	}
	cols, err := rows.Columns()
	if err != nil {
		return
	}
	fmt.Println("Rows ", cols)
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

func NewClient(driver string, config config) (*GormConnection, error) {
	init, ok := supportedDrivers[driver]
	if !ok {
		return nil, driverNotSupported
	}
	// opt := {}
	return newGormConnection(init(config))
}

type driverInitializer = func(config) driver

var supportedDrivers = map[string]driverInitializer{
	"postgres": newPostgresDriver,
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

func NewConfig(username, password, dbname string) config {
	return config{
		Host:     "localhost",
		Username: username,
		DBName:   dbname,
		Port:     "5432",
		Password: password,
	}
}

type QueryResult struct {
	rows        *sql.Rows
	processTime time.Duration
	err         error
}
