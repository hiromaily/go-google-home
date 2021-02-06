package commands

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/hiromaily/go-google-home/pkg/device"
	"go.uber.org/zap"
	"os"
)

// base command

func newWrapperCmd(
	logger *zap.Logger,
	devicer device.Device,
	chFinishNotifier chan struct{},
	cmd subcommands.Command,
) *wrapperCmd {
	return &wrapperCmd{
		logger:           logger,
		devicer:          devicer,
		chFinishNotifier: chFinishNotifier,
		Command:          cmd,
	}
}

type wrapperCmd struct {
	logger           *zap.Logger
	devicer          device.Device
	chFinishNotifier chan struct{}

	subcommands.Command
	help bool
}

func (w *wrapperCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&w.help, "help", false, "show help")
	w.Command.SetFlags(f)
}

func (w *wrapperCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// help
	if w.help {
		f.Parse([]string{os.Args[1]})
		return subcommands.HelpCommand().Execute(ctx, f, args...)
	}

	// devicer initialization
	_, err := w.devicer.Start()
	if err != nil {
		w.logger.Error("fail to call devicer.Start()", zap.Error(err))
		return subcommands.ExitFailure
	}

	defer func() {
		w.devicer.Close()
		close(w.chFinishNotifier)
	}()

	// execute
	exitStatus := w.Command.Execute(ctx, f, args...)
	w.logger.Debug("exitStatus", zap.Int("exitStatus", int(exitStatus)))

	// wait
	<-w.chFinishNotifier
	return exitStatus
}
