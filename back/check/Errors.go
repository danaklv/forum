package check

import (
	"fr/back/models"
	"html/template"
	"net/http"
)

func Status500(w http.ResponseWriter) {
	w.WriteHeader(500)
	tmpl, err := template.ParseFiles("./front/templates/500.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := models.ErrData{Result: "Internal Server Error"}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func Status400(w http.ResponseWriter) {
	w.WriteHeader(400)
	tmpl, err := template.ParseFiles("./front/templates/400.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := models.ErrData{Result: "Bad Request"}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func Status404(w http.ResponseWriter) {
	w.WriteHeader(404)
	tmpl, err := template.ParseFiles("./front/templates/404.html")
	if err != nil {
		Status500(w)
		return
	}
	data := models.ErrData{Result: "Not Found"}
	if err := tmpl.Execute(w, data); err != nil {
		Status500(w)
		return
	}
}

func Status401(w http.ResponseWriter) {
	w.WriteHeader(401)
	tmpl, err := template.ParseFiles("./front/templates/notauthenticated.html")
	if err != nil {
		Status500(w)
		return
	}
	data := models.ErrData{Result: "Unauthorized"}
	if err := tmpl.Execute(w, data); err != nil {
		Status500(w)
		return
	}

}

func Status405(w http.ResponseWriter) {
	w.WriteHeader(405)
	tmpl, err := template.ParseFiles("./front/templates/405.html")
	if err != nil {
		Status500(w)
		return
	}
	data := models.ErrData{Result: "Method not allowed"}
	if err := tmpl.Execute(w, data); err != nil {
		Status500(w)
		return
	}
}