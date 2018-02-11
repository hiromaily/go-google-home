package googlehome

// Google Text to Speech API

import (
	"context"
	"fmt"
	"github.com/barnybug/go-cast"
	ctl "github.com/barnybug/go-cast/controllers"
	lg "github.com/hiromaily/golibs/log"
	"github.com/micro/mdns"
	"net"
	ur "net/url"
	"strings"
	"time"
)

const (
	castService = "_googlecast._tcp"
	ttsURL      = "https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob"
	ghPrefix    = "Google-Home-"
)

type GoogleHome struct {
	host   string
	AddrV4 net.IP
	Port   int
	Error  error
	Controller
}

type Controller struct {
	client *cast.Client
	ctx    context.Context
}

func New(strIP string, port int) *GoogleHome {
	ip := net.ParseIP(strIP)
	gh := GoogleHome{AddrV4: ip, Port: port}
	return &gh
}

func DiscoverService() *GoogleHome {
	notifyService := make(chan *GoogleHome)

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 1)

	var isDone bool
	go func() {
		//	for entry := range entriesCh {
		//	}
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

				gh := GoogleHome{Error: fmt.Errorf("Timeout for discovering devices.")}
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

func (g *GoogleHome) NewClient() error {
	ctx := context.Background()
	client := cast.NewClient(g.AddrV4, g.Port)
	err := client.Connect(ctx)
	if err != nil {
		return err
	}

	lg.Infof("Connected to %v:%d", g.AddrV4, g.Port)
	g.Controller = Controller{client: client, ctx: ctx}
	return nil
}

func (c *Controller) Speak(text string, language string) error {
	u := fmt.Sprintf(ttsURL, ur.QueryEscape(text), ur.QueryEscape(language))
	return c.Play(u)
}

func (c *Controller) Play(url string) error {
	media, err := c.client.Media(c.ctx)
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

func (c *Controller) Stop() error {
	if !c.client.IsPlaying(c.ctx) {
		return nil
	}
	media, err := c.client.Media(c.ctx)
	if err != nil {
		return err
	}
	_, err = media.Stop(c.ctx)
	return err
}

func (c *Controller) Quit() error {
	receiver := c.client.Receiver()
	_, err := receiver.QuitApp(c.ctx)
	return err
}

func (c *Controller) Close() {
	c.client.Close()
}
