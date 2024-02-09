package clarum

import (
	"github.com/goclarum/clarum/core/config"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/core/logging"
)

func Setup() {
	logging.Infof("Starting clarum %s", config.Version())
}

func Finish() {
	logging.Info("Waiting for all actions to finish.")

	control.RunningActions.Wait()

	logging.Info("All actions finished.")
}
