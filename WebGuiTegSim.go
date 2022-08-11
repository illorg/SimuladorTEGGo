package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type DetallesConsulta struct {
	FichasAtaque  string
	FichasDefensa string
	Simulaciones  string
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

		// do something with details
		//_ = details

		tmpl.Execute(w, struct{ Success bool }{true})
		fmt.Println(details)
		a, _ := strconv.Atoi(details.FichasAtaque)
		b, _ := strconv.Atoi(details.FichasDefensa)
		s, _ := strconv.Atoi(details.Simulaciones)

		Maestro(a, b, s)

	})

	http.ListenAndServe(":8080", nil)
}
