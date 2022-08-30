package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var victoriasataque int
var victoriasdefensa int
var porcentajevictoria float32
var tiempo string

type DetallesConsulta struct {
	FichasAtaque  string
	FichasDefensa string
	Simulaciones  string
}
type ResultadosSim struct {
	Success           bool
	VictoriasLabel    int
	DerrotasLabel     int
	PorcentajeLabel   float32
	SimulacionesLabel int
	TiempoLabel       string
}

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := DetallesConsulta{
			FichasAtaque:  r.FormValue("fichasataque"),
			FichasDefensa: r.FormValue("fichasdefensa"),
			Simulaciones:  r.FormValue("simulaciones"),
		}

		a, _ := strconv.Atoi(details.FichasAtaque)
		b, _ := strconv.Atoi(details.FichasDefensa)
		s, _ := strconv.Atoi(details.Simulaciones)

		victoriasataque, victoriasdefensa, porcentajevictoria, tiempo = SimularTeg(a, b, s)
		resultadosHtml := ResultadosSim{Success: true, VictoriasLabel: victoriasataque, DerrotasLabel: victoriasdefensa,
			PorcentajeLabel: porcentajevictoria, SimulacionesLabel: s, TiempoLabel: tiempo}
		tmpl.Execute(w, resultadosHtml)
		fmt.Println(details)

	})

	http.ListenAndServe(":8080", nil)
}
