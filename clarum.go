package clarum

import (
	"fmt"
	"github.com/goclarum/clarum/core/config"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http"
	"log/slog"
)

// Entry point for HTTP endpoints configuration
func Http() http.Builder {
	return &http.EndpointBuilder{}
}

func Setup() {
	slog.Info(fmt.Sprintf("Starting clarum %s", config.Version()))
}

func Finish() {
	slog.Info("Waiting for all actions to finish.")

	control.RunningActions.Wait()

	slog.Info("All actions finished.")
}
