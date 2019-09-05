package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jimsmart/schema"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
	"github.com/smallnest/gen/dbmeta"
	gtmpl "github.com/smallnest/gen/template"
)

var (
	sqlType     = "mysql"
	sqlConnStr  = "root:123456@tcp(127.0.0.1:3306)/drfinder?"
	sqlTable    = "posts"

	packageName = ""

	jsonAnnotation = false
	gormAnnotation = true
	gureguTypes    = false
	rest = false
)

func main() {

	//driverSource := "root:123456@tcp(127.0.0.1:3306)/drfinder?"
	var db, err = sql.Open(sqlType, sqlConnStr)
	if err != nil {
		fmt.Println("Error in open database: " + err.Error())
		return
	}
	defer db.Close()

	// parse or read tables
	var tables []string
	if sqlTable != "" {
		tables = strings.Split(sqlTable, ",")
	} else {
		tables, err = schema.TableNames(db)
		if err != nil {
			fmt.Println("Error in fetching tables information from mysql information schema")
			return
		}
	}

	os.Mkdir("model", 0777)

	apiName := "api"
	if rest {
		os.Mkdir(apiName, 0777)
	}

	t, err := getTemplate(gtmpl.ModelTmpl)
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return
	}

	ct, err := getTemplate(gtmpl.ControllerTmpl)
	if err != nil {
		fmt.Println("Error in loading controller template: " + err.Error())
		return
	}

	var structNames []string

	// generate go files for each table
	for _, tableName := range tables {
		structName := dbmeta.FmtFieldName(tableName)
		structName = inflection.Singular(structName)
		structNames = append(structNames, structName)

		modelInfo := dbmeta.GenerateStruct(db, tableName, structName, "model", jsonAnnotation, gormAnnotation, gureguTypes)

		var buf bytes.Buffer
		err = t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering model: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join("model", inflection.Singular(tableName)+".go"), data, 0777)

		if rest {
			//write api
			buf.Reset()
			err = ct.Execute(&buf, map[string]string{"PackageName": packageName + "/model", "StructName": structName})
			if err != nil {
				fmt.Println("Error in rendering controller: " + err.Error())
				return
			}
			data, err = format.Source(buf.Bytes())
			if err != nil {
				fmt.Println("Error in formating source: " + err.Error())
				return
			}
			ioutil.WriteFile(filepath.Join(apiName, inflection.Singular(tableName)+".go"), data, 0777)
		}
	}

	if rest {
		rt, err := getTemplate(gtmpl.RouterTmpl)
		if err != nil {
			fmt.Println("Error in lading router template")
			return
		}
		var buf bytes.Buffer
		err = rt.Execute(&buf, structNames)
		if err != nil {
			fmt.Println("Error in rendering router: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(apiName, "router.go"), data, 0777)
	}
}

func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"pluralize":        inflection.Plural,
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func camelToLowerCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])

	return strings.Join(ss, "")
}