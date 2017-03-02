package handlers

import (
	"iris.arke.works/forum/http/helper"
	"net/http"
)

func DenyHandler(w http.ResponseWriter, r *http.Request) {
	helper.ErrorStringWriter(w, r, http.StatusMethodNotAllowed, "Method not allowed")
}
