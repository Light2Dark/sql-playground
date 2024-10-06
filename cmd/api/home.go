package main

import (
	"fmt"
	"net/http"

	"github.com/Light2Dark/sql-playground/internal/templates"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffcontext"
)

func (app application) homeHandler(w http.ResponseWriter, r *http.Request) {
	user := ffcontext.NewEvaluationContext("unique-key-2")
	hasFlag, err := ffclient.BoolVariation("test-flag", user, false)
	if err != nil {
		app.logger.Error(fmt.Sprintf("Error with feature flag evaluation: %s", err))
	}

	if hasFlag {
		app.logger.Info("got stuf")
	} else {
		app.logger.Info("disabled!")
	}

	templates.Home().Render(r.Context(), w)
}
