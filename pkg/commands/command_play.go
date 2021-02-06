package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/device"
)

// playCmd defines play command
type playCmd struct {
	logger           *zap.Logger
	devicer          device.Device
	chFinishNotifier chan struct{}

	// args
	musicURL string
}

func newPlayCmd(
	logger *zap.Logger,
	devicer device.Device,
	chFinishNotifier chan struct{},
) *playCmd {
	return &playCmd{
		logger:           logger,
		devicer:          devicer,
		chFinishNotifier: chFinishNotifier,
	}
}

func (*playCmd) Name() string {
	return "play"
}

func (*playCmd) Synopsis() string {
	return "play sound data"
}

func (c *playCmd) Usage() string {
	return fmt.Sprintf(`Usage: gh play [options]... : %s
 options:
  -url    url of sound data like mp3
 e.g. 
  gh play -url "http://xxxxx/music/play.mp"
`, c.Synopsis())
}

func (c *playCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.musicURL, "url", "", "url of sound data")
}

func (c *playCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.musicURL == "" {
		fmt.Println(c.Usage())
		return subcommands.ExitUsageError
	}
	c.logger.Info("plays", zap.String("url", c.musicURL))

	// set callback event
	c.devicer.Controller().RunEventReceiver(c.chFinishNotifier)

	// speak message
	err := c.devicer.Controller().Play(c.musicURL)
	if err != nil {
		c.logger.Error("fail to call Speak()", zap.Error(err))
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
