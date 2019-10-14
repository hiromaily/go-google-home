package googlehome

// Google Text to Speech API
// https://www.w3.org/TR/speech-synthesis/

import (
	"context"
	"fmt"

	"github.com/barnybug/go-cast"
	ctl "github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"

	"net"
	ur "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bookerzzz/grok"
	"github.com/micro/mdns"

	//ev "github.com/barnybug/go-cast/events"
	"github.com/pkg/errors"

	lg "github.com/hiromaily/golibs/log"
)

const (
	castService = "_googlecast._tcp"
	ttsURL      = "https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob"
	ghPrefix    = "Google-Home-"
)

// GoogleHome is GoogleHome object
type GoogleHome struct {
	host   string
	AddrV4 net.IP
	Port   int
	Error  error
	Controller
}

// Controller is for controlling google home by cast.Client
type Controller struct {
	Client *cast.Client
	ctx    context.Context
}

// NewGoogleHome is to return empty GoogleHome object
func NewGoogleHome() *GoogleHome {
	return &GoogleHome{}
}

// NewGoogleHomeWithAddress is to return GoogleHome object
func NewGoogleHomeWithAddress(strIP string, port int) *GoogleHome {
	ip := net.ParseIP(strIP)
	gh := GoogleHome{AddrV4: ip, Port: port}
	return &gh
}

// WithAddressString is ad address, port to GoogleHome object
func (g *GoogleHome) WithAddressString(address string) (*GoogleHome, error) {
	//use address if it exists.
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		return nil, errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
	}
	g.AddrV4 = net.ParseIP(addr[0])
	g.Port = port
	return g, nil
}

// DiscoverService is to discover google home devices
func DiscoverService() *GoogleHome {
	notifyService := make(chan *GoogleHome)

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 1)

	var isDone bool
	go func() {
		for {
			select {
			case entry := <-entriesCh:
				lg.Info("Discovered Device.")
				if isDone {
					return
				}
				lg.Debugf("Name: %s", entry.Name)
				lg.Debugf("Host: %s", entry.Host)
				lg.Debugf("AddrV4: %v", entry.AddrV4)
				lg.Debugf("Port: %d", entry.Port)

				//e.g. Name: Google-Home-1234567890abcdefghijklmn._googlecast._tcp.local.
				if strings.HasPrefix(entry.Name, ghPrefix) {
					isDone = true
					gh := GoogleHome{host: entry.Host, AddrV4: entry.AddrV4, Port: entry.Port}
					notifyService <- &gh
					close(entriesCh)
					return
				}
			case <-time.After(5 * time.Second):
				isDone = true
				close(entriesCh)

				gh := GoogleHome{Error: fmt.Errorf("timeout for discovering devices")}
				notifyService <- &gh
				return
			}
		}
	}()

	// Start the lookup
	mdnsLookup(entriesCh)

	// receiver for waiting
	return <-notifyService
}

func mdnsLookup(entriesCh chan *mdns.ServiceEntry) {
	params := mdns.DefaultParams(castService)
	params.Entries = entriesCh
	//params.WantUnicastResponse = true
	mdns.Query(params)
}

// NewClient is to create client for google cast controller
func (g *GoogleHome) NewClient() error {
	ctx := context.Background()
	client := cast.NewClient(g.AddrV4, g.Port)
	err := client.Connect(ctx)
	if err != nil {
		return err
	}

	lg.Infof("Connected to %v:%d", g.AddrV4, g.Port)
	g.Controller = Controller{Client: client, ctx: ctx}
	return nil
}

// Speak is to speak by text
func (c *Controller) Speak(text string, language string) error {
	u := fmt.Sprintf(ttsURL, ur.QueryEscape(text), ur.QueryEscape(language))
	return c.Play(u)
}

// Play is to play music by url
func (c *Controller) Play(url string) error {
	media, err := c.Client.Media(c.ctx)
	if err != nil {
		return err
	}

	item := ctl.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: "audio/mpeg",
	}
	_, err = media.LoadMedia(c.ctx, item, 0, true, map[string]interface{}{})
	return err
}

// Stop is stop playing music
func (c *Controller) Stop() error {
	if !c.Client.IsPlaying(c.ctx) {
		return nil
	}
	media, err := c.Client.Media(c.ctx)
	if err != nil {
		return err
	}
	_, err = media.Stop(c.ctx)
	return err
}

// GetStatus is to get google cast client status
func (c *Controller) GetStatus() (*ctl.MediaStatusResponse, error) {
	media, err := c.Client.Media(c.ctx)
	if err != nil {
		return nil, err
	}

	//*MediaStatusResponse, error
	return media.GetStatus(c.ctx)
}

// SetVolume is to set volume
func (c *Controller) SetVolume(vol string) error {
	receiver := c.Client.Receiver()
	level, _ := strconv.ParseFloat(vol, 64)
	muted := false
	volume := ctl.Volume{Level: &level, Muted: &muted}
	_, err := receiver.SetVolume(c.ctx, &volume)
	if err != nil {
		return err
	}
	return nil
}

// Close is to close google cast client
func (c *Controller) Close() {
	c.Client.Close()
}

// RunEventReceiver is to receive current event
// FIXME: it seems to be useless.
func (c *Controller) RunEventReceiver(notify chan bool) {
	go func() {
		for evt := range c.Client.Events {
			//TODO:evt is type of interface, it should be casted to something.
			lg.Infof("[Event received] %v", evt)
			//switch evt.(type)
			//S1034: assigning the result of this type assertion to a variable (switch evt := evt.(type)) could eliminate
			switch evt := evt.(type) {
			case ctl.MediaStatus:
				//grok.Value(evt)
				lg.Debugf("PlayerState: %s", evt.PlayerState)
				if evt.IdleReason == "FINISHED" {
					lg.Debug("ctl.MediaStatus: FINISHED")
					notify <- true
				}
			case events.AppStarted:
				lg.Debug("AppStarted")
			case events.AppStopped:
				lg.Debug("AppStopped")
			case events.Connected:
				lg.Debug("Connected")
			case events.Disconnected:
				lg.Debug("Disconnected")
			case events.StatusUpdated:
				lg.Debug("StatusUpdated")
			default:
				lg.Debug("events: default")
				grok.Value(evt)
			}
		}
	}()
}

// DebugStatus is for debugging of status
func (c *Controller) DebugStatus(status *ctl.MediaStatusResponse) {
	fmt.Println("DebugStatus(): *ctl.MediaStatusResponse:status")
	grok.Value(status)
}
