package device

import (
	"context"
	"net"
	"strconv"
	"strings"

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
	WithIPPort(ip string, port int) Device
	WithAddress(address string) Device
	WithService(srv *Service) Device
	WithClient() Device
	Controller() controller.Controller
	Error() error
}

// Device object
type device struct {
	logger          *zap.Logger
	serviceReceiver ServiceReceiver
	ctl             controller.Controller
	host            string
	addrV4          net.IP
	port            int
	err             error
}

// NewGoogleHome is to return empty GoogleHome object
func NewDevice(logger *zap.Logger, serviceReceiver ServiceReceiver) Device {
	return &device{
		logger:          logger,
		serviceReceiver: serviceReceiver,
	}
}

// Start
func (d *device) Start(addr string) (Device, error) {
	if addr != "" {
		// create object from address
		d.logger.Debug("Start()", zap.String("address", addr))
		return d.WithAddress(addr).WithClient(), nil
	}
	// discover Google Home
	d.logger.Debug("discover google home address")
	srv := d.serviceReceiver.Discover()
	if srv.Error != nil {
		return d.WithService(srv).WithClient(), nil
	}
	return nil, srv.Error
}

// NewGoogleHomeWithAddress is to return GoogleHome object
func (d *device) WithIPPort(ip string, port int) Device {
	parsedIP := net.ParseIP(ip)
	d.addrV4 = parsedIP
	d.port = port
	return d
}

// WithAddressString is ad address, port to GoogleHome object
func (d *device) WithAddress(address string) Device {
	if d.err != nil {
		return d
	}
	// use address if it exists.
	//	host, strPort, err := net.SplitHostPort(u.Host)
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		d.err = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return d
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		d.err = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return d
	}
	return d.WithIPPort(addr[0], port)
}

// WithClient is to create client for google cast controller
func (d *device) WithClient() Device {
	if d.err != nil {
		return d
	}
	// create client
	ctx := context.Background()
	client := cast.NewClient(d.addrV4, d.port)
	err := client.Connect(ctx)
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

// WithClient is to create client for google cast controller
func (d *device) WithService(srv *Service) Device {
	d.host = srv.Service.Host
	d.addrV4 = srv.Service.AddrV4
	d.port = srv.Service.Port
	return d
}

func (d *device) Controller() controller.Controller {
	return d.ctl
}

func (d *device) Error() error {
	return d.err
}
