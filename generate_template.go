package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	strcase "github.com/stoewer/go-strcase"
)

// generateSQL crea el archivo sql
func generateSQL(m Model, d string) {
	now := time.Now()
	fn := "sql_" + now.Format("20060102") + "_" + now.Format("150405") + "_create_" + m.Table + ".sql"
	generateTemplate(filepath.Join(d, fn), "table_sql.gotpl", m)
}

// generateSQL crea el archivo sql
func generatePSQL(m Model, d string) {
	now := time.Now()
	fn := "psql_" + now.Format("20060102") + "_" + now.Format("150405") + "_create_" + m.Table + ".sql"
	generateTemplate(filepath.Join(d, fn), "table_psql.gotpl", m)
}

// generateSQL crea el archivo sql
func generateOSQL(m Model, d string) {
	now := time.Now()
	fn := "osql_" + now.Format("20060102") + "_" + now.Format("150405") + "_create_" + m.Table + ".sql"
	generateTemplate(filepath.Join(d, fn), "table_osql.gotpl", m)
}

// generateApplication crea el Application
func generateApplication(m Model, d string) {
	generateTemplate(filepath.Join(d, "application_service.go"), "application_service.gotpl", m)
}

// generateDomain crea el domain
func generateDomain(m Model, d string) {
	generateTemplate(filepath.Join(d, "domain.go"), "domain.gotpl", m)
}

// generatePorts crea la interface ports
func generateStorage(m Model, d string) {
	generateTemplate(filepath.Join(d, "storage.go"), "storage.gotpl", m)
}

// generateSqlServer crea el archivo sqlserver
func generateSqlServer(m Model, d string) {
	generateTemplate(filepath.Join(d, "repository_sqlserver.go"), "repository_sqlserver.gotpl", m)
}

// generatePsql crea el archivo psql
func generatePsql(m Model, d string) {
	generateTemplate(filepath.Join(d, "repository_psql.go"), "repository_psql.gotpl", m)
}

// generateOracle crea el archivo oracle
func generateOracle(m Model, d string) {
	generateTemplate(filepath.Join(d, "repository_oracle.go"), "repository_oracle.gotpl", m)
}

// generateApplication crea el Application
func generateHandler(m Model, d string) {
	generateTemplate(filepath.Join(d, "handler.go"), "handler.gotpl", m)
}

// generateApplication crea el Application
func generateRouter(m Model, d string) {
	generateTemplate(filepath.Join(d, "router.go"), "router.gotpl", m)
}

// generateApplication crea el Application
func generateModel(m Model, d string) {
	generateTemplate(filepath.Join(d, "model.go"), "model.gotpl", m)
}

// generateTemplate crea el archivo .go con base al template
func generateTemplate(dest, source string, m Model) {
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("no se pudo crear el archivo: %v", err)
	}
	if filepath.Ext(dest) == ".go" {
		defer formatFile(dest)
	}
	defer f.Close()

	var tpl *template.Template
	switch m.ID {
	case "int":
		tpl = tplInt
	case "int64":
		tpl = tplInt64
	default:
		tpl = tplUuid
	}

	err = tpl.ExecuteTemplate(f, source, m)
	if err != nil {
		log.Printf("error creando el archivo: %v", err)
		return
	}
}

func setHelpers() {
	fm = template.FuncMap{
		"ucc": func(v string) string {
			return strcase.UpperCamelCase(v)
		},
		"upp": func(v string) string {
			return strings.ToUpper(v)
		},
		"lwc": func(v string) string {
			return strings.ToLower(v)
		},
		"kcc": func(v string) string {
			return strcase.KebabCase(v)
		},
		"lcc": func(v string) string {
			return strcase.LowerCamelCase(v)
		},
		"inc": func(v int) int {
			return v + 1
		},
		"dec": func(v int) int {
			return v - 1
		},
		"pkg": func(v string) string {
			return strings.Split(v, ".")[1]
		},
		"mdlType": func(v string) string {
			switch v {
			case "uuid":
				return "string"
			default:
				return v
			}
			return "CHANGE-THIS-TYPE"
		},
		"sqlType": func(v string) string {
			switch v {
			case "uint":
				fallthrough
			case "int":
				return "INT"
			case "int64":
				return "BIGINT"
			case "string":
				return "VARCHAR"
			case "bool":
				return "BOOLEAN"
			case "time.Time":
				return "TIMESTAMP"
			case "uuid":
				return "UNIQUEIDENTIFIER"
			}
			return "CHANGE-THIS-TYPE"
		},
		"psqlType": func(v string) string {
			switch v {
			case "uint":
				fallthrough
			case "int":
				return "INTEGER"
			case "int64":
				return "BIGINT"
			case "string":
				return "VARCHAR"
			case "bool":
				return "BOOLEAN"
			case "time.Time":
				return "TIMESTAMP"
			case "uuid":
				return "UUID"
			}
			return "CHANGE-THIS-TYPE"
		},
		"osqlType": func(v string) string {
			switch v {
			case "uint":
				fallthrough
			case "int":
				return "NUMBER(10,0)"
			case "int64":
				return "NUMBER(20,0)"
			case "string":
				return "VARCHAR2"
			case "bool":
				return "NUMBER(1)"
			case "time.Time":
				return "TIMESTAMP"
			case "uuid":
				return "VARCHAR2(50)"
			}
			return "CHANGE-THIS-TYPE"
		},
		"fieldSQL": func(f Field) string {
			field := strcase.UpperCamelCase(f.Name)

			if f.NotNull == "NOT NULL" {
				return fmt.Sprintf("m.%s", field)
			}

			switch f.Type {
			case "string":
				return fmt.Sprintf("psql.StringToNull(m.%s)", field)
			case "int":
				fallthrough
			case "uint":
				return fmt.Sprintf("psql.IntToNull(int64(m.%s))", field)
			case "time.Time":
				return fmt.Sprintf("psql.TimeToNull(m.%s)", field)
			default:
				return fmt.Sprintf("Error: no existe el tipo de dato: %s", t)
			}
		},
		"fieldSQLScan": func(f Field) string {
			if f.NotNull == "NOT NULL" {
				return ""
			}

			switch f.Type {
			case "string":
				return fmt.Sprintf("%s := sql.NullString{}", f.Name)
			case "int":
				fallthrough
			case "uint":
				return fmt.Sprintf("%s := sql.NullInt64{}", f.Name)
			case "time.Time":
				return fmt.Sprintf("%s := pq.NullTime{}", f.Name)
			case "bool":
				return fmt.Sprintf("%s := sql.NullBool{}", f.Name)
			default:
				return fmt.Sprintf("Error: no existe el tipo de dato: %s", t)
			}
		},
		"fieldSQLScanValue": func(f Field) string {
			field := strcase.UpperCamelCase(f.Name)
			if f.NotNull == "NOT NULL" {
				return ""
			}
			switch f.Type {
			case "string":
				return fmt.Sprintf("m.%s = %s.String", field, f.Name)
			case "int":
				fallthrough
			case "uint":
				return fmt.Sprintf("m.%s = %s(%s.Int64)", field, f.Type, f.Name)
			case "time.Time":
				return fmt.Sprintf("m.%s = %s.Time", field, f.Name)
			case "bool":
				return fmt.Sprintf("m.%s = %s.Bool", field, f.Name)
			default:
				return fmt.Sprintf("Error: no existe el tipo de dato: %s", t)
			}
		},
	}
}
