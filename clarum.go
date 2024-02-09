package clarum

import (
	"github.com/go-clarum/clarum-core/config"
	"github.com/go-clarum/clarum-core/control"
	"github.com/go-clarum/clarum-core/logging"
)

func Setup() {
	logging.Infof("Starting clarum %s", config.Version())
}

func Finish() {
	logging.Info("Waiting for all actions to finish.")

	control.RunningActions.Wait()

	logging.Info("All actions finished.")
}
