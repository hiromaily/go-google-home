package commands

import (
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/device"
)

// Register registers sum commands
func Register(logger *zap.Logger, devicer device.Device, chFinishNotifier chan struct{}) {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	// subcommands.Register(newSpeakCmd(logger, devicer, chFinishNotifier), "speak")
	subcommands.Register(
		newWrapperCmd(logger, devicer, chFinishNotifier, newSpeakCmd(logger, devicer, chFinishNotifier)),
		"speak",
	)
	// subcommands.Register(newPlayCmd(logger, devicer, chFinishNotifier), "play")
	subcommands.Register(
		newWrapperCmd(logger, devicer, chFinishNotifier, newPlayCmd(logger, devicer, chFinishNotifier)),
		"play",
	)

	// subcommands.Register(newServerCmd(logger, devicer), "server")

	flag.Parse()
}
