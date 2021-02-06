package device

import (
	"strings"
	"time"

	"github.com/micro/mdns"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//-----------------------------------------------------------------------------
// ServiceServiceReceiver
//-----------------------------------------------------------------------------

type ServiceReceiver interface {
	Discover() *Service
}

type serviceReceiver struct {
	logger *zap.Logger
}

func NewServiceReceiver(logger *zap.Logger) ServiceReceiver {
	return &serviceReceiver{
		logger: logger,
	}
}

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
				s.logger.Info("Discovered Device.")
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
			case <-time.After(5 * time.Second):
				isDone = true
				close(chEntry)
				chNotifier <- &Service{Error: errors.Errorf("timeout for discovering devices")}
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
