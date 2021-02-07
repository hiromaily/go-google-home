package main

import (
	"time"

	"github.com/hiromaily/go-google-home/pkg/server"

	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/config"
	"github.com/hiromaily/go-google-home/pkg/device"
	"github.com/hiromaily/go-google-home/pkg/logger"
)

// Registry interface
type Registry interface {
	NewDevicer() device.Device
	NewServer() server.Server
	NewLogger() *zap.Logger
}

type registry struct {
	conf   *config.Root
	server device.Device
	logger *zap.Logger
}

// NewRegistry return Registry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewDevicer return device.Device interface
func (r *registry) NewDevicer() device.Device {
	if r.server == nil {
		r.server = device.NewDevice(
			r.NewLogger(),
			r.newServiceReceiver(),
			r.conf.Device.Address,
			r.conf.Device.Lang,
		)
	}
	return r.server
}

// NewDevicer returns device.Device interface
func (r *registry) NewServer() server.Server {
	return server.NewServer(
		r.NewLogger(),
		r.NewDevicer(),
		r.conf.Server.Port,
	)
}

func (r *registry) NewLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newServiceReceiver() device.ServiceReceiver {
	parsedDuration, err := time.ParseDuration(r.conf.Device.Timeout)
	if err != nil {
		panic(err)
	}
	return device.NewServiceReceiver(r.NewLogger(), parsedDuration)
}
