package fileshandlers

import (
	"net/http"
)

func ScriptHandler(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./front/js/script.js")
}
