package main

import (
	"go.uber.org/zap"
	"time"

	"github.com/hiromaily/go-google-home/pkg/config"
	"github.com/hiromaily/go-google-home/pkg/device"
	"github.com/hiromaily/go-google-home/pkg/logger"
)

// Registry interface
type Registry interface {
	NewDevicer() device.Device
	NewLogger() *zap.Logger
}

type registry struct {
	conf   *config.Root
	logger *zap.Logger
}

// NewRegistry return Registry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewDevicer return device.Device interface
func (r *registry) NewDevicer() device.Device {
	return device.NewDevice(
		r.NewLogger(),
		r.newServiceReceiver(),
		r.conf.Device.Address,
		r.conf.Device.Lang,
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
