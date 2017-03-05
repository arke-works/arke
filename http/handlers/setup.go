package handlers

import (
	"github.com/pressly/chi"
	"net/http"
)

func MakeRouter(r chi.Router) {
	r.Get("/:resource/:snowflake", GetHandler)
	r.Get("/:resource", GetHandler)
	r.Head("/:resource/:snowflake", GetHandler)
	r.Head("/:resource", GetHandler)
	r.Options("/:resource", OptionHandler)
	r.Options("/:resource/:unused", OptionHandler)
	r.Post("/:resource", PostHandler)
	r.Post("/:resource/:unused", DenyHandler)
	r.Delete("/:resource/", DenyHandler)
	r.Delete("/:resource/:snowflake", DeleteHandler)
	r.Put("/:resource", DenyHandler)
	r.Put("/:resource/:snowflake", DenyHandler)
	r.Patch("/:resource", DenyHandler)
	r.Patch("/:resource/:snowflake", DenyHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Route not found"))
	})
}
