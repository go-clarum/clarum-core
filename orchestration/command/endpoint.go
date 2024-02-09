package command

import (
	"context"
	"fmt"
	"github.com/go-clarum/clarum-core/config"
	"github.com/go-clarum/clarum-core/control"
	"github.com/go-clarum/clarum-core/durations"
	"github.com/go-clarum/clarum-core/logging"
	clarumstrings "github.com/go-clarum/clarum-core/validators/strings"
	"os/exec"
	"time"
)

type Endpoint struct {
	name          string
	cmdComponents []string
	warmup        time.Duration
	cmd           *exec.Cmd
	cmdCancel     context.CancelFunc
	logger        *logging.Logger
}

func newCommandEndpoint(name string, components []string, warmup time.Duration) *Endpoint {
	if len(components) == 0 || clarumstrings.IsBlank(components[0]) {
		panic("Builder cannot start anything - cmd is empty")
	}

	return &Endpoint{
		name:          name,
		cmdComponents: components,
		warmup:        durations.GetDurationWithDefault(warmup, 1*time.Second),
		logger:        logging.NewLogger(config.LoggingLevel(), logPrefix(name)),
	}
}

// Start the process from the given command & arguments.
// The process will be started into a cancelable context so that we can
// cancel it later in the post-integration test phase.
func (endpoint *Endpoint) start() error {
	endpoint.logger.Infof("running cmd [%s]", endpoint.cmdComponents)
	ctx, cancel := context.WithCancel(context.Background())

	endpoint.cmd = exec.CommandContext(ctx, endpoint.cmdComponents[0], endpoint.cmdComponents[1:]...)
	endpoint.cmdCancel = cancel

	endpoint.logger.Debug("starting command")
	if err := endpoint.cmd.Start(); err != nil {
		return err
	} else {
		endpoint.logger.Debug("cmd start successful")
	}

	time.Sleep(endpoint.warmup)
	endpoint.logger.Debug("warmup ended")

	return nil
}

// Stop the running process. Since the process was created with a context, we will attempt to
// call ctx.Cancel(). If it returns an error, the process will be killed just in case.
// We also wait for the action here, so that the post-integration test phase ends successfully.
func (endpoint *Endpoint) stop() error {
	control.RunningActions.Add(1)
	defer control.RunningActions.Done()

	endpoint.logger.Infof("stopping cmd [%s]", endpoint.cmdComponents)

	if endpoint.cmdCancel != nil {
		endpoint.logger.Debug("cancelling cmd")
		endpoint.cmdCancel()

		if _, err := endpoint.cmd.Process.Wait(); err != nil {
			endpoint.logger.Errorf("cmd.Wait() returned error - [%s]", err)
			endpoint.killProcess()
			return err
		} else {
			endpoint.logger.Debug("context cancel finished successfully")
		}
	} else {
		if err := endpoint.cmd.Process.Release(); err != nil {
			endpoint.logger.Errorf("cmd.Release() returned error - [%s]", err)
			endpoint.killProcess()
			return err
		} else {
			endpoint.logger.Debug("cmd kill successful")
		}
	}

	return nil
}

func (endpoint *Endpoint) killProcess() {
	endpoint.logger.Info("killing process")

	if err := endpoint.cmd.Process.Kill(); err != nil {
		endpoint.logger.Errorf("cmd.Kill() returned error - [%s]", err)
		return
	}
}

func logPrefix(cmdName string) string {
	return fmt.Sprintf("Command %s: ", cmdName)
}
