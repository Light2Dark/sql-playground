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

	data := r.Form.Get("editor")
	app.logger.Debug("submit called", "data", data)

	rows, err := app.db.Query(r.Context(), `SELECT * from users`)
	if err != nil {
		app.logger.Error("error querying the database", "error", err)
		templates.Message(fmt.Sprintf("Error: %v", err)).Render(r.Context(), w)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		panic(err)
	}

	templates.Message(fmt.Sprintf("Results: %v", res)).Render(r.Context(), w)
}
