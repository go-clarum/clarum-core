package command

import (
	"context"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/core/durations"
	clarumstrings "github.com/goclarum/clarum/core/validators/strings"
	"log/slog"
	"os/exec"
	"time"
)

type Endpoint struct {
	name          string
	cmdComponents []string
	warmup        time.Duration
	cmd           *exec.Cmd
	cmdCancel     context.CancelFunc
}

func newCommandEndpoint(name string, components []string, warmup time.Duration) *Endpoint {
	if len(components) == 0 || clarumstrings.IsBlank(components[0]) {
		panic("Builder cannot start anything - cmd is empty")
	}

	return &Endpoint{
		name:          name,
		cmdComponents: components,
		warmup:        durations.GetDurationWithDefault(warmup, 1*time.Second),
	}
}

// Start the process from the given command & arguments.
// The process will be started into a cancelable context so that we can
// cancel it later in the post-integration test phase.
func (endpoint *Endpoint) start() error {
	logPrefix := logPrefix(endpoint.name)
	slog.Info(fmt.Sprintf("%s: running cmd [%s]", logPrefix, endpoint.cmdComponents))
	ctx, cancel := context.WithCancel(context.Background())

	endpoint.cmd = exec.CommandContext(ctx, endpoint.cmdComponents[0], endpoint.cmdComponents[1:]...)
	endpoint.cmdCancel = cancel

	slog.Debug(fmt.Sprintf("%s: starting command", logPrefix))
	if err := endpoint.cmd.Start(); err != nil {
		return err
	} else {
		slog.Debug(fmt.Sprintf("%s: cmd start successful", logPrefix))
	}

	time.Sleep(endpoint.warmup)
	slog.Debug(fmt.Sprintf("%s: warmup ended", logPrefix))

	return nil
}

// Stop the running process. Since the process was created with a context, we will attempt to
// call ctx.Cancel(). If it returns an error, the process will be killed just in case.
// We also wait for the action here, so that the post-integration test phase ends successfully.
func (endpoint *Endpoint) stop() error {
	control.RunningActions.Add(1)
	defer control.RunningActions.Done()

	logPrefix := logPrefix(endpoint.name)
	slog.Info(fmt.Sprintf("%s: stopping cmd [%s]", logPrefix, endpoint.cmdComponents))

	if endpoint.cmdCancel != nil {
		slog.Debug(fmt.Sprintf("%s: cancelling cmd", logPrefix))
		endpoint.cmdCancel()

		if _, err := endpoint.cmd.Process.Wait(); err != nil {
			slog.Error(fmt.Sprintf(fmt.Sprintf("%s: cmd.Wait() returned error - [%s]", logPrefix, err)))
			endpoint.killProcess()
			return err
		} else {
			slog.Debug(fmt.Sprintf("%s: context cancel finished successfully", logPrefix))
		}
	} else {
		if err := endpoint.cmd.Process.Release(); err != nil {
			slog.Error(fmt.Sprintf("%s: cmd.Release() returned error - [%s]", logPrefix, err))
			endpoint.killProcess()
			return err
		} else {
			slog.Debug(fmt.Sprintf("%s: cmd kill successful", logPrefix))
		}
	}

	return nil
}

func (endpoint *Endpoint) killProcess() {
	logPrefix := logPrefix(endpoint.name)
	slog.Info(fmt.Sprintf(fmt.Sprintf("%s: killing process", logPrefix)))

	if err := endpoint.cmd.Process.Kill(); err != nil {
		slog.Error(fmt.Sprintf("%s: cmd.Kill() returned error - [%s]", logPrefix, err))
		return
	}
}

func logPrefix(cmdName string) string {
	return fmt.Sprintf("Command %s", cmdName)
}
