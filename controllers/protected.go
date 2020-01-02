package controllers

import (
	"database/sql"
	"../utils"
	"net/http"
)

func (c Controller) Protected(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "Invoked Endpoint")
	}
}