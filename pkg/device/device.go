package device

import (
	"context"
	"net"
	"strconv"

	"github.com/barnybug/go-cast"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/controller"
)

const (
	castService = "_googlecast._tcp"
	ghPrefix    = "Google-Home-"
)

//-----------------------------------------------------------------------------
// Device
//-----------------------------------------------------------------------------

// Device interface
type Device interface {
	Start(addr string) (Device, error)
	withIPPort(ip string, port int) Device
	WithAddress(address string) Device
	WithService(srv *Service) Device
	WithClient() Device
	Controller() controller.Controller
	Error() error
	Close()
}

type device struct {
	logger          *zap.Logger
	serviceReceiver ServiceReceiver
	ctl             controller.Controller
	host            string
	addrV4          net.IP
	port            int
	err             error
}

// NewDevice returns Device interface
func NewDevice(logger *zap.Logger, serviceReceiver ServiceReceiver) Device {
	return &device{
		logger:          logger,
		serviceReceiver: serviceReceiver,
	}
}

// Start starts setup for Device
func (d *device) Start(addr string) (Device, error) {
	if addr != "" {
		// create object from address
		d.logger.Debug("Start()", zap.String("address", addr))
		return d.WithAddress(addr).WithClient(), nil
	}
	// discover Google Home
	d.logger.Debug("discover google home address")
	srv := d.serviceReceiver.Discover()
	if srv.Error == nil {
		return d.WithService(srv).WithClient(), nil
	}
	return nil, srv.Error
}

// withIPPort sets ip, port to device object
func (d *device) withIPPort(ip string, port int) Device {
	parsedIP := net.ParseIP(ip)
	d.addrV4 = parsedIP
	d.port = port
	return d
}

// WithAddress sets address to device object
func (d *device) WithAddress(address string) Device {
	if d.err != nil {
		return d
	}
	host, strPort, err := net.SplitHostPort(address)
	if err != nil {
		d.err = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return d

	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		d.err = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return d
	}
	return d.withIPPort(host, port)
}

// WithClient is to create client for google cast controller
func (d *device) WithClient() Device {
	if d.err != nil {
		return d
	}
	// create client
	ctx := context.Background()
	client, err := d.connect(ctx)
	if err != nil {
		d.err = errors.Errorf("fail to connect by google cast")
		return d
	}
	d.logger.Info("connected",
		zap.Any("address", d.addrV4),
		zap.Int("port", d.port),
	)

	// create controller
	d.ctl = controller.NewController(ctx, client, d.logger)
	return d
}

func (d *device) connect(ctx context.Context) (*cast.Client, error) {
	client := cast.NewClient(d.addrV4, d.port)
	err := client.Connect(ctx)
	if err != nil {
		return nil, errors.Errorf("fail to connect by google cast")
	}
	return client, nil
}

// WithClient sets host,address, port from Service to device object
func (d *device) WithService(srv *Service) Device {
	d.host = srv.Service.Host
	d.addrV4 = srv.Service.AddrV4
	d.port = srv.Service.Port
	return d
}

// Controller returns controller.Controller
func (d *device) Controller() controller.Controller {
	return d.ctl
}

// Error returns error
func (d *device) Error() error {
	return d.err
}

// Close closes device
func (d *device) Close() {
	d.ctl.Close()
}
