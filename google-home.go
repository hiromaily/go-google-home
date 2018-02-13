package googlehome

// Google Text to Speech API
// https://www.w3.org/TR/speech-synthesis/

import (
	"context"
	"fmt"
	"github.com/barnybug/go-cast"
	ctl "github.com/barnybug/go-cast/controllers"
	ev "github.com/barnybug/go-cast/events"
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
	Client *cast.Client
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
	g.Controller = Controller{Client: client, ctx: ctx}
	return nil
}

func (c *Controller) Speak(text string, language string) error {
	u := fmt.Sprintf(ttsURL, ur.QueryEscape(text), ur.QueryEscape(language))
	return c.Play(u)
}

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

func (c *Controller) GetStatus() (*ctl.MediaStatusResponse, error) {
	media, err := c.Client.Media(c.ctx)
	if err != nil {
		return nil, err
	}

	//*MediaStatusResponse, error
	return media.GetStatus(c.ctx)
}

func (c *Controller) Close() {
	c.Client.Close()
}

func (c *Controller) RunEventReceiver() {
	go func() {
		for evt := range c.Client.Events {
			//TODO:evt is type of interface, it should be casted to something.
			fmt.Println("[Event received]", evt)
			switch evt.(type) {
			case ctl.MediaStatus:
				fmt.Println("ctl.MediaStatus")
				if obj, ok := evt.(ctl.MediaStatus); ok {
					fmt.Println(obj)
				}
			case ev.Connected:
				fmt.Println("ev.Connected")
				if obj, ok := evt.(ev.Connected); ok {
					fmt.Println(obj)
				}
			case ev.Disconnected:
				fmt.Println("ev.Disconnected")
				if obj, ok := evt.(ev.Disconnected); ok {
					fmt.Println(obj)
				}
			case ev.StatusUpdated:
				fmt.Println("ev.StatusUpdated")
				if obj, ok := evt.(ev.StatusUpdated); ok {
					fmt.Println(obj)
				}
			//case ctl.MediaStatusResponse:
			//	fmt.Println("OK!!")
			//case ctl.ReceiverStatus:
			//	fmt.Println("OK2!!")
			//case ctl.StatusResponse:
			//	fmt.Println("OK3!!")
			default:
				fmt.Println("default")
				//TODO: what is this event type??
				//{CC1AD845 Default Media Receiver Ready To Cast}
			}
		}
	}()
}

func (c *Controller) DebugStatus(status *ctl.MediaStatusResponse) {
	//type MediaStatus struct {
	//	net.PayloadHeaders
	//	MediaSessionID         int                    `json:"mediaSessionId"`
	//	PlaybackRate           float64                `json:"playbackRate"`
	//	PlayerState            string                 `json:"playerState"`
	//	CurrentTime            float64                `json:"currentTime"`
	//	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	//	Volume                 *Volume                `json:"volume,omitempty"`
	//	Media                  *MediaStatusMedia      `json:"media"`
	//	CustomData             map[string]interface{} `json:"customData"`
	//	RepeatMode             string                 `json:"repeatMode"`
	//	IdleReason             string                 `json:"idleReason"`
	//}

	fmt.Println("[status:length]", len(status.Status)) //1
	fmt.Printf("[status[0]:obj] %#v\n", status.Status[0])
	fmt.Println("[status[0].Type]", status.Status[0].Type)                           // none
	fmt.Println("[status[0].RequestId]", status.Status[0].RequestId)                 // <nil>
	fmt.Println("[status[0].CurrentTime]", status.Status[0].CurrentTime)             // 0
	fmt.Println("[status[0].CustomData]", status.Status[0].CustomData)               // map[]
	fmt.Println("[status[0].IdleReason]", status.Status[0].IdleReason)               // none
	fmt.Println("[status[0].Media.ContentId]", status.Status[0].Media.ContentId)     // url for playing
	fmt.Println("[status[0].Media.ContentType]", status.Status[0].Media.ContentType) // audio/mpeg
	fmt.Println("[status[0].Media.Duration]", status.Status[0].Media.Duration)       // 0.96
	fmt.Println("[status[0].Media.StreamType]", status.Status[0].Media.StreamType)   // BUFFERED

	fmt.Println("[*status[0].Volume.Level]", *status.Status[0].Volume.Level)   // 1
	fmt.Println("[status[0].PayloadHeaders]", status.Status[0].PayloadHeaders) //{ <nil>}

	fmt.Println("[*status.RequestId]", *status.RequestId)                   //0xc42007ed68
	fmt.Println("[status.Type]", status.Type)                               //MEDIA_STATUS
	fmt.Println("[status.PayloadHeaders.Type]", status.PayloadHeaders.Type) //MEDIA_STATUS
}
