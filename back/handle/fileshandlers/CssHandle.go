package fileshandlers

import (
	"net/http"
	"os"
)

func CssHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./front/css/style.css"); os.IsNotExist(err) {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "./front/css/style.css")
}
