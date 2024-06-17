package fileshandlers

import "net/http"

func CheckAuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./front/templates/notauthenticated.html")
}
