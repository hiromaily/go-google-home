package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/hiromaily/go-google-home/pkg/server"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

// serverCmd defines server command
type serverCmd struct {
	logger *zap.Logger
	server server.Server

	// args
	port int
}

func newServerCmd(
	logger *zap.Logger,
	server server.Server,
) *serverCmd {
	return &serverCmd{
		logger: logger,
		server: server,
	}
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "run as server"
}

func (c *serverCmd) Usage() string {
	return fmt.Sprintf(`Usage: gh server [options]... : %s
 options:
  -port    server port
 e.g. 
  gh server -port 8888
`, c.Synopsis())
}

func (c *serverCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.port, "port", 0, "server port")
}

func (c *serverCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c.logger.Info("server", zap.Int("port", c.port))

	// speak message
	err := c.server.Start(c.port)
	if err != nil {
		c.logger.Error("fail to call Start()", zap.Error(err))
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
