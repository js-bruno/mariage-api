package utils

import (
	"log/slog"
	"os"
)

// TODO: create a trace logic to follow-up the code flow
func SetStructuredLogging() {
	logging := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logging)
}
