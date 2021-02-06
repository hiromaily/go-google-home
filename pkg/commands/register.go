package commands

import (
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/device"
)

func Register(logger *zap.Logger, devicer device.Device, chFinishNotifier chan struct{}) {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(newSpeakCmd(logger, devicer, chFinishNotifier), "msg")
	//subcommands.Register(newMusicCmd(devicer), "music")
	//subcommands.Register(newServerCmd(devicer), "server")

	flag.Parse()
}

//func validateArguments() {
//	// this pattern is not allowed
//	if !*server && *message == "" && *music == "" {
//		flag.Usage()
//		os.Exit(1)
//		return
//	}
//}
