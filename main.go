package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

var (
	tplUuid  *template.Template
	tplInt   *template.Template
	tplInt64 *template.Template
	fm       = template.FuncMap{}
	// id tipo de dato para crear el ID
	l  string
	id string
	h  string
	// n nombre del paqute
	n string
	// t nombre de la tabla
	t string
	// fs los campos del modelo
	fs []Field
	// rutas de los paquetes de configuración, logger, message, model_role
	ps map[string]string
)

func init() {
	setHelpers()
	tplUuid = template.Must(template.New("").Funcs(fm).ParseGlob("templates/uuid/*.gotpl"))
	tplInt = template.Must(template.New("").Funcs(fm).ParseGlob("templates/int/*.gotpl"))
	tplInt64 = template.Must(template.New("").Funcs(fm).ParseGlob("templates/int64/*.gotpl"))
}

func main() {

	readConfigFile("./config.json")
	ff := ps["src"]

	color.Green("Iniciando proceso...")

	if ff != "" {
		generateFromFile(ff)
	} else {
		showMainMenu()
		execute()
	}

	color.Green("Proceso finalizado.")
}

func execute() {
	m := Model{id, n, t, fs, ps, l}
	gopath := os.Getenv("GOPATH")
	realDest := []string{gopath, "src"}
	realDest = append(realDest, strings.Split(ps["dest"], "/")...)
	gp := filepath.Join(realDest...)
	schema := strings.Split(t, ".")
	if len(schema) == 2 {
		t = fmt.Sprintf("%s/%s", schema[0], schema[1])
	}
	pks := filepath.Join(gp, t)
	ds := filepath.Join(gp, "_database/sqlserver")
	dsp := filepath.Join(gp, "_database/postgres")
	dso := filepath.Join(gp, "_database/oracle")

	createDir(pks)
	createDir(ds)
	createDir(dsp)
	createDir(dso)
	generateApplication(m, pks)
	generateDomain(m, pks)
	generateStorage(m, pks)
	generateSqlServer(m, pks)
	generatePsql(m, pks)
	generateOracle(m, pks)
	generateSQL(m, ds)
	generatePSQL(m, dsp)
	generateOSQL(m, dso)
	// generador de apis
	if h == "api" {
		api := filepath.Join(gp, "handler/"+schema[1])
		createDir(api)
		generateHandler(m, api)
		generateRouter(m, api)
		generateModel(m, api)
	}
}

// createDir crea el directorio de destino de los archivos
func createDir(d string) {
	_, err := os.Stat(d)
	if os.IsNotExist(err) {
		log.Printf("no existe la carpeta %s. Creandola...", d)
		os.MkdirAll(d, os.ModePerm)
	}
}

func formatFile(filePath string) {
	cmd := exec.Command("gofmt", "-w", filePath)
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: No se pudo ejecutar gofmt")
	}
}

func readConfigFile(cf string) {
	ps = make(map[string]string)

	file, err := ioutil.ReadFile(cf)
	if err != nil {
		e := fmt.Sprintf("no se pudo abrir el archivo de configuración: %v", err)
		color.Red(e)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &ps)
	if err != nil {
		e := fmt.Sprintf("no se pudo convertir la configuración en mapa: %v", err)
		color.Red(e)
		os.Exit(1)
	}
}

func getFields(value string) []Field {
	var err error
	rs := make([]Field, 0)
	fields := strings.Split(value, " ")
	for _, v := range fields {
		field := strings.Split(v, ":")
		nn := "NOT NULL"
		i := 0
		if len(field) >= 3 {
			if strings.ToLower(field[2]) == "t" {
				nn = ""
			}
		}
		if len(field) == 4 {
			i, err = strconv.Atoi(field[3])
			if err != nil {
				log.Fatalf("%s no es un número válido: %v", field[3], err)
			}

		}
		f := Field{field[0], field[1], nn, i}
		rs = append(rs, f)
	}

	return rs
}
