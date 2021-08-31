package dbms

import (
	"database/sql"
	"errors"
	"time"
)

var supportedDrivers = map[string]func(Config) Driver{
	"postgres": NewPostgres,
}

type Driver interface {
	Connect(config Config) error
	Query(text string, args ...interface{}) (*sql.Rows, error)
	Close()
}
type DBMS interface {
	Connect(config Config) error
	Close()
}
type dbms struct {
	driver Driver
}

var (
	DriverNotSupported = errors.New("Driver Not Supportd")
)

func New(driver string, config Config) (DBMS, error) {

	_func, ok := supportedDrivers[driver]
	if !ok {
		return nil, DriverNotSupported
	}
	return dbms{
		driver: _func(config),
	}, nil
}
func (sys *dbms) Query(text string, args ...interface{}) QueryResult {
	start := time.Now()
	sys.driver.Query(text, args...)
	qr := QueryResult{}
	qr.processTime = time.Now().Sub(start)
	return qr
}
func (sys dbms) Connect(config Config) error {
	return sys.driver.Connect(config)
}
func (sys dbms) Close() {
	sys.driver.Close()
}

type Config struct {
	username, host, port, password string
}
type QueryResult struct {
	rows        *sql.Rows
	processTime time.Duration
	err         error
}
