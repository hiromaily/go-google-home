package main

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/config"
	"github.com/hiromaily/go-google-home/pkg/device"
	"github.com/hiromaily/go-google-home/pkg/logger"
)

// Registry interface
type Registry interface {
	NewDevicer() device.Device
}

type registry struct {
	conf   *config.Root
	logger *zap.Logger
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewBooker is to register for booker interface
func (r *registry) NewDevicer() device.Device {
	return device.NewDevice(
		r.newLogger(),
		r.newServiceReceiver(),
	)
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newServiceReceiver() device.ServiceReceiver {
	return device.NewServiceReceiver(r.newLogger())
}
