package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

func generateFromFile(file string) {
	isCSV := strings.HasSuffix(file, ".csv")
	if !isCSV {
		color.Red("el archivo de importación de paquetes debe tener extensión .csv")
		os.Exit(1)
	}

	f, err := os.Open(file)
	if err != nil {
		color.Red(fmt.Sprintf("no se pudo abrir el archivo de importación de paquetes: %v", err))
		os.Exit(1)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			color.Red(fmt.Sprintf("error leyendo la línea del archivo de importación de paquetes: %v", err))
			os.Exit(1)
		}

		l = record[0]
		h = record[1]
		id = record[2]
		n = record[3]
		t = record[4]
		fields := record[5]

		if t == "" {
			color.Red(fmt.Sprintf("no se procesó el modelo: %s porque no se recibieron campos", id))
			continue
		}

		fs = getFields(fields)

		execute()
	}
}
