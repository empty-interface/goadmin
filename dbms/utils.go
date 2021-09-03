package dbms

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Table struct {
	Catalog                      string `gorm:"column:table_catalog"`
	Schema                       string `gorm:"column:table_schema"`
	Name                         string `gorm:"column:table_name"`
	Type                         string `gorm:"column:table_type"`
	Self_referencing_column_name string `gorm:"column:self_referencing_column_name"`
}

func (*Table) TableName() string {

	return "information_schema.tables"
}

func (conn *GormConnection) GetTables() []Table {
	// query := `select * from information_schema.tables where table_schema = 'public'`
	// tables := make([]Table, 0)
	// rows, err := conn.Raw(query).Rows()
	// if err != nil {
	// 	return
	// }
	// rows.Scan()
	tables := []Table{}
	conn.Where("table_schema=?", "public").Find(&tables)
	fmt.Println("Found", len(tables), "tables")
	return tables
}

type TableInfo struct {
	Types []*sql.ColumnType
}

func (conn *GormConnection) GetTableRows(name string, limit int) ([][]interface{}, []string, error) {
	query := fmt.Sprintf(`select * from %s limit %v`, name, limit)
	rows, err := conn.Raw(query).Rows()
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	var result [][]interface{}
	row := make([]interface{}, len(cols))
	temp := make([]interface{}, len(cols))
	for i := range temp {
		temp[i] = &row[i]
	}
	for rows.Next() {
		rows.Scan(temp...)
		result = append(result, row)
	}
	// fmt.Println("--------------------------------\n ")
	// fmt.Println("data", result)
	// fmt.Println("\n--------------------------------")
	return result, cols, nil
}
func (conn *GormConnection) GetTableInfo(name string) (*TableInfo, error) {
	query := fmt.Sprintf(`select * from %s`, name)
	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	return &TableInfo{
		Types: types,
	}, nil
}
