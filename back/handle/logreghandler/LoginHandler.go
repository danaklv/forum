package logreghandler

import (
	"fmt"
	"fr/back/check"
	"html/template"
	"net/http"
)

func LoginHadnler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("front/templates/reg.html", "front/templates/header.html")
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}
	tmpl.ExecuteTemplate(w, "reg", nil)
}
