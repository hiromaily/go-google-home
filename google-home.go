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
)

const (
	castService = "_googlecast._tcp"
	ttsURL      = "https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob"
	ghPrefix    = "Google-Home-"
)

//type ServiceEntry struct {
//	Name       string
//	Host       string
//	AddrV4     net.IP
//	AddrV6     net.IP
//	Port       int
//	Info       string
//	InfoFields []string
//	TTL        int
//
//	Addr net.IP // @Deprecated
//
//	hasTXT bool
//	sent   bool
//}

type GoogleHome struct {
	host   string
	addrV4 net.IP
	port   int
	Controller
}

type Controller struct {
	client *cast.Client
	ctx    context.Context
}

func DiscoverService() *GoogleHome {
	//TODO:timeout and return error
	notifyService := make(chan *GoogleHome)

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			lg.Info("Discovered Device.")
			lg.Debugf("Name: %s", entry.Name)
			lg.Debugf("Host: %s", entry.Host)
			lg.Debugf("AddrV4: %v", entry.AddrV4)
			lg.Debugf("Port: %d", entry.Port)
			//lg.Debugf("Info: %s", entry.Info)
			//lg.Debugf("InfoFields: %v", entry.InfoFields)
			//lg.Debugf("TTL: %v", entry.TTL)

			//e.g. Name: Google-Home-1234567890abcdefghijklmn._googlecast._tcp.local.
			if strings.Contains(entry.Name, ghPrefix) {
				gh := GoogleHome{host: entry.Host, addrV4: entry.AddrV4, port: entry.Port}
				notifyService <- &gh

				close(entriesCh)
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
	client := cast.NewClient(g.addrV4, g.port)
	err := client.Connect(ctx)
	if err != nil {
		return err
	}

	lg.Infof("Connected to %v:%d", g.addrV4, g.port)
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
