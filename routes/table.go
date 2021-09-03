package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/empty-interface/goadmin/dbms"
)

const TablePath = "/table"

var sections = []string{"structure", "select", "add", "custom"}

func HandleTable(w http.ResponseWriter, r *http.Request, currentSession *Session) (int, error) {
	tableName := getTableName(r)
	if tableName == "" {
		return http.StatusBadRequest, fmt.Errorf("Invalid table name")
	}
	tmpl := template.Must(template.ParseFiles("html/table.html"))
	sectionToShow := getSection(r, "structure")
	page := tablePage{
		Showing:   sectionToShow,
		TableName: tableName,
	}

	switch sectionToShow {
	case "structure":
		str, err := showStructure(currentSession, tableName)
		if err != nil {
			return http.StatusBadRequest, err
		}
		page.Structure = str
	case "select":
		page.Select = &selectSection{}
		query := getQuery(r, tableName, 50)
		rows, names, err := showSelect(currentSession, query, 50)

		page.Select.Query = query
		page.Select.Error = false
		if err != nil {
			page.Select.Error = true
			page.Select.ErrMsg = err.Error()
		} else {
			fmt.Println("Data length", len(rows))
			page.Select.Rows = rows
			page.Select.Names = names
		}
	case "custom":

	}
	tmpl.Execute(w, &page)
	return -1, nil

}

func showSelect(currentSession *Session, name string, limit int) ([][]interface{}, []string, error) {
	return currentSession.Conn.GetTableRows(name, limit)
}
func showStructure(currentSession *Session, name string) (*dbms.TableInfo, error) {
	tableInfo, err := currentSession.Conn.GetTableInfo(name)
	if err != nil {
		return nil, err
	}
	return tableInfo, err
}
func getTableName(r *http.Request) string {
	name, exist := r.URL.Query()["name"]
	if !exist {
		return ""
	}
	return name[0]
}

func getSection(r *http.Request, defaultSection string) string {
	section, exist := r.URL.Query()["section"]
	if !exist || !validSection(section[0]) {
		return defaultSection
	}
	return section[0]
}

func validSection(section string) bool {
	for _, sec := range sections {
		if sec == section {
			return true
		}
	}
	return false
}
func getQuery(r *http.Request, tableName string, limit int) string {
	queries, exist := r.URL.Query()["query"]
	if !exist || queries[0] == "" {
		return fmt.Sprintf("select * from %s limit %v", tableName, limit)
	}
	return queries[0]
}

type tablePage struct {
	Showing   string
	TableName string
	Structure *dbms.TableInfo
	Select    *selectSection
	// Add       *dbms.TableInfo
}
type selectSection struct {
	Rows   [][]interface{}
	Names  []string
	Query  string
	Error  bool
	ErrMsg string
}
