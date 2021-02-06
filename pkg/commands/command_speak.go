package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/device"
)

// speakCmd defines args
type speakCmd struct {
	logger           *zap.Logger
	devicer          device.Device
	chFinishNotifier chan struct{}
	lang             string

	// args
	message string
}

func newSpeakCmd(
	logger *zap.Logger,
	devicer device.Device,
	chFinishNotifier chan struct{},
) *speakCmd {
	return &speakCmd{
		logger:           logger,
		devicer:          devicer,
		chFinishNotifier: chFinishNotifier,
	}
}

func (*speakCmd) Name() string {
	return "speak"
}

func (*speakCmd) Synopsis() string {
	return "speak message"
}

func (*speakCmd) Usage() string {
	return `usage: msg`
}

func (c *speakCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.message, "msg", "", "message")
}

func (c *speakCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.message == "" {
		c.logger.Warn("-msg is required")
		return subcommands.ExitFailure
	}
	c.logger.Info("speaks", zap.String("msg", c.message))

	// set callback event
	c.devicer.Controller().RunEventReceiver(c.chFinishNotifier)

	// speak message
	err := c.devicer.Controller().Speak(c.message, c.devicer.Lang())
	if err != nil {
		c.logger.Error("fail to call Speak()", zap.Error(err))
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
