package googlehome

// Google Text to Speech API
// https://www.w3.org/TR/speech-synthesis/

import (
	"context"
	"fmt"
	"net"
	ur "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/barnybug/go-cast"
	ctl "github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
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

//-----------------------------------------------------------------------------
// Controller
//-----------------------------------------------------------------------------

// Controller is Controller interface
type Controller interface {
	Speak(text string, language string) error
	Play(url string) error
	Stop() error
	GetStatus() (*ctl.MediaStatusResponse, error)
	SetVolume(vol string) error
	Close()
	RunEventReceiver(notify chan bool)
	DebugStatus(status *ctl.MediaStatusResponse)
}

//-----------------------------------------------------------------------------
// GoogleHome
//-----------------------------------------------------------------------------

// GoogleHome is GoogleHome object
type GoogleHome struct {
	host   string
	AddrV4 net.IP
	Port   int
	Error  error
	//ctl Controller
	Control
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
func (g *GoogleHome) WithAddressString(address string) *GoogleHome {
	if g.Error != nil {
		return g
	}
	//use address if it exists.
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		g.Error = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return g
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		g.Error = errors.Errorf("address is invalid. format should be :%s", "xxx.xxx.xxx.xxx:8009")
		return g
	}
	g.AddrV4 = net.ParseIP(addr[0])
	g.Port = port
	return g
}

// WithClient is to create client for google cast controller
func (g *GoogleHome) WithClient() *GoogleHome {
	if g.Error != nil {
		return g
	}
	ctx := context.Background()
	client := cast.NewClient(g.AddrV4, g.Port)
	err := client.Connect(ctx)
	if err != nil {
		g.Error = errors.Errorf("fail to connect by google cast")
		return g
	}

	lg.Infof("Connected to %v:%d", g.AddrV4, g.Port)
	//g.ctl = &Control{Client: client, ctx: ctx}
	g.Control = Control{Client: client, ctx: ctx}
	return g
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
				// set error
				gh := GoogleHome{Error: errors.Errorf("timeout for discovering devices")}
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

//-----------------------------------------------------------------------------
// Controller
//-----------------------------------------------------------------------------

// Control is for controlling google home by cast.Client
type Control struct {
	Client *cast.Client
	ctx    context.Context
}

// Speak is to speak by text
func (c *Control) Speak(text string, language string) error {
	u := fmt.Sprintf(ttsURL, ur.QueryEscape(text), ur.QueryEscape(language))
	return c.Play(u)
}

// Play is to play music by url
func (c *Control) Play(url string) error {
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
func (c *Control) Stop() error {
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
func (c *Control) GetStatus() (*ctl.MediaStatusResponse, error) {
	media, err := c.Client.Media(c.ctx)
	if err != nil {
		return nil, err
	}

	//*MediaStatusResponse, error
	return media.GetStatus(c.ctx)
}

// SetVolume is to set volume
func (c *Control) SetVolume(vol string) error {
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
func (c *Control) Close() {
	c.Client.Close()
}

// RunEventReceiver is to receive current event
func (c *Control) RunEventReceiver(notify chan bool) {
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
func (c *Control) DebugStatus(status *ctl.MediaStatusResponse) {
	fmt.Println("DebugStatus(): *ctl.MediaStatusResponse:status")
	grok.Value(status)
}
