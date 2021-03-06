package device

import (
	"strings"
	"time"

	//"github.com/micro/mdns"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/mdns"
)

// TODO: maybe better
// https://github.com/hashicorp/mdns

//-----------------------------------------------------------------------------
// ServiceServiceReceiver
//-----------------------------------------------------------------------------

// ServiceReceiver interface
type ServiceReceiver interface {
	Discover() *Service
}

type serviceReceiver struct {
	logger  *zap.Logger
	timeout time.Duration
}

// NewServiceReceiver returns ServiceReceiver
func NewServiceReceiver(logger *zap.Logger, timeout time.Duration) ServiceReceiver {
	return &serviceReceiver{
		logger:  logger,
		timeout: timeout,
	}
}

// Service includes *mdns.ServiceEntry
type Service struct {
	Service *mdns.ServiceEntry
	Error   error
}

// Discover discovers google home devices
func (s *serviceReceiver) Discover() *Service {
	chNotifier := make(chan *Service)
	chEntry := make(chan *mdns.ServiceEntry, 1)

	var isDone bool
	go func() {
		for {
			select {
			case entry := <-chEntry:
				if isDone {
					return
				}
				s.logger.Info("Discovered Device.",
					zap.String("name", entry.Name),
					zap.String("host", entry.Host),
					zap.Any("addr_v4", entry.AddrV4),
					zap.Int("port", entry.Port),
				)

				// e.g. Name: Google-Home-1234567890abcdefghijklmn._googlecast._tcp.local.
				if strings.HasPrefix(entry.Name, ghPrefix) {
					isDone = true
					chNotifier <- &Service{Service: entry}
					close(chEntry)
					return
				}
			case <-time.After(s.timeout):
				isDone = true
				close(chEntry)
				chNotifier <- &Service{Error: errors.Errorf("can't discover devices by timeout")}
				return
			}
		}
	}()

	// start lookup
	mdnsLookup(chEntry)

	// receiver for waiting
	return <-chNotifier
}

func mdnsLookup(chEntry chan *mdns.ServiceEntry) {
	params := mdns.DefaultParams(castService)
	params.Entries = chEntry
	// params.WantUnicastResponse = true
	mdns.Query(params)
}
