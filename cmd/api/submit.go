package main

import (
	"fmt"
	"net/http"

	"github.com/Light2Dark/sql-playground/internal/templates"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (app application) submitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	query := r.Form.Get("editor")
	app.logger.Debug("submit called", "query", query)

	rows, err := app.db.Query(r.Context(), query)
	if err != nil {
		templates.Message(fmt.Sprintf("Error: %v", err)).Render(r.Context(), w)
		return
	}
	defer rows.Close()

	var res [][]any
	for rows.Next() {
		rowResponse, err := rows.Values()
		if err != nil {
			app.logger.Error(fmt.Sprintf("Error: %v"))
		}
		res = append(res, rowResponse)
	}

	templates.Message(fmt.Sprintf("Results: %v", res)).Render(r.Context(), w)
}
