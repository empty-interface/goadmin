package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/empty-interface/goadmin/dbms"
)

const TablePath = "/table"

var sections = []string{"structure", "select", "add"}

func HandleTable(w http.ResponseWriter, r *http.Request, currentSession *Session) (int, error) {
	name := getTableName(r)
	if name == "" {
		return http.StatusBadRequest, fmt.Errorf("Invalid table name")
	}
	tmpl := template.Must(template.ParseFiles("html/table.html"))
	sectionToShow := getSection(r, "structure")
	page := tablePage{
		Showing:   sectionToShow,
		TableName: name,
	}

	switch sectionToShow {
	case "structure":
		str, err := showStructure(currentSession, name)
		if err != nil {
			return http.StatusBadRequest, err
		}
		page.Structure = str
	case "select":
		rows, names, err := showSelect(currentSession, name, 50)
		if err != nil {
			return http.StatusBadRequest, err
		}
		fmt.Println("Data length", len(rows))
		page.Select = &selectSection{
			Rows:  rows,
			Names: names,
		}
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

type tablePage struct {
	Showing   string
	TableName string
	Structure *dbms.TableInfo
	Select    *selectSection
	// Add       *dbms.TableInfo
}
type selectSection struct {
	Rows  [][]interface{}
	Names []string
}
