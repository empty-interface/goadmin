package dbms

import "fmt"

type Table struct {
	Table_catalog                string `gorm:"column:table_catalog"`
	Table_schema                 string `gorm:"column:table_schema"`
	Table_name                   string `gorm:"column:table_name"`
	Table_type                   string `gorm:"column:table_type"`
	Self_referencing_column_name string `gorm:"column:self_referencing_column_name"`
}

func (*Table) TableName() string {

	return "information_schema.tables"
}

func (conn *GormConnection) GetTables() {
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
}
