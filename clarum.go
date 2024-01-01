package clarum

import (
	"fmt"
	"github.com/goclarum/clarum/core/config"
	"github.com/goclarum/clarum/core/control"
	clarumhttp "github.com/goclarum/clarum/http"
	"log/slog"
)

// Entry point for HTTP endpoints configuration
func Http() *clarumhttp.EndpointBuilder {
	return &clarumhttp.EndpointBuilder{}
}

func Setup() {
	slog.Info(fmt.Sprintf("Starting clarum %s", config.Version()))

	// TODO: go 1.22 will allow us to set the level on the default logger without changing the format
	//  this change is what we need: https://github.com/golang/go/commit/3188758653fc7d2b229e234273d41878ddfdd5f2
	//  release date February 2024
	//h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: config.LoggingLevel()})
	//slog.SetDefault(slog.New(h))
}

func Finish() {
	slog.Info("Waiting for all actions to finish.")

	control.RunningActions.Wait()

	slog.Info("All actions finished.")
}
