package controller

// https://www.w3.org/TR/speech-synthesis/

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	ttsURL = "https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob"
)

//-----------------------------------------------------------------------------
// Controller
//-----------------------------------------------------------------------------

// Controller interface
type Controller interface {
	Speak(text string, language string) error
	Play(url string) error
	Stop() error
	GetStatus() (*controllers.MediaStatusResponse, error)
	SetVolume(vol string) error
	Close()
	CloseEvent()
	RunEventReceiver(notify chan bool)
	DebugStatus(status *controllers.MediaStatusResponse)
}

// controller controls google home by cast.Client
type controller struct {
	ctx    context.Context
	client *cast.Client
	logger *zap.Logger
}

// NewController returns Controller
func NewController(
	ctx context.Context,
	client *cast.Client,
	logger *zap.Logger,
) Controller {
	return &controller{
		ctx:    ctx,
		client: client,
		logger: logger,
	}
}

// Speak speaks text content
func (c *controller) Speak(text string, language string) error {
	playURL := fmt.Sprintf(ttsURL, url.QueryEscape(text), url.QueryEscape(language))
	c.logger.Debug("speak", zap.String("playURL", playURL))
	return c.Play(playURL)
}

// Play plays music of URL
func (c *controller) Play(playURL string) error {
	media, err := c.client.Media(c.ctx)
	if err != nil {
		return errors.Wrap(err, "fail to call cast.Client.Media()")
	}

	mediaItem := controllers.MediaItem{
		ContentId:   playURL,
		StreamType:  "BUFFERED",
		ContentType: "audio/mpeg",
	}
	customData := map[string]interface{}{}

	_, err = media.LoadMedia(c.ctx, mediaItem, 0, true, customData)
	return err
}

// Stop stops playing music
func (c *controller) Stop() error {
	if !c.client.IsPlaying(c.ctx) {
		return nil
	}
	media, err := c.client.Media(c.ctx)
	if err != nil {
		return errors.Wrap(err, "fail to call cast.Client.Media()")
	}
	_, err = media.Stop(c.ctx)
	return err
}

// GetStatus returns media status
func (c *controller) GetStatus() (*controllers.MediaStatusResponse, error) {
	media, err := c.client.Media(c.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call cast.Client.Media()")
	}

	//*MediaStatusResponse, error
	return media.GetStatus(c.ctx)
}

// SetVolume sets volume
func (c *controller) SetVolume(vol string) error {
	receiver := c.client.Receiver()
	level, _ := strconv.ParseFloat(vol, 64)
	muted := false
	volume := controllers.Volume{
		Level: &level,
		Muted: &muted,
	}
	_, err := receiver.SetVolume(c.ctx, &volume)
	if err != nil {
		return errors.Wrap(err, "fail to call cast.Client.Receiver.SetVolume()")
	}
	return nil
}

// Close closes google cast client
func (c *controller) Close() {
	c.client.Close()
}

// CloseEvent closes google cast client Events channel
func (c *controller) CloseEvent() {
	close(c.client.Events)
}

// RunEventReceiver runs event receiver of client status
func (c *controller) RunEventReceiver(notify chan bool) {
	go func() {
		for evt := range c.client.Events {
			// evt is type of interface, it should be asserted
			c.logger.Info("event_received",
				zap.Any("events.Event", evt),
			)
			switch evtType := evt.(type) {
			case controllers.MediaStatus:
				c.logger.Debug("controllers.MediaStatus",
					zap.String("evtType.PlayerState", evtType.PlayerState),
					zap.String("evtType.IdleReason", evtType.IdleReason),
				)
				if evtType.IdleReason == "FINISHED" {
					c.logger.Debug("controllers.MediaStatus: FINISHED")
					notify <- true
				}
			case events.AppStarted:
				c.logger.Debug("evtType: AppStarted")
			case events.AppStopped:
				c.logger.Debug("evtType: AppStopped")
			case events.Connected:
				c.logger.Debug("evtType: Connected")
			case events.Disconnected:
				c.logger.Debug("evtType: Disconnected")
			case events.StatusUpdated:
				c.logger.Debug("evtType: StatusUpdated")
			default:
				c.logger.Debug("evtType: default")
				grok.Value(evt)
			}
		}
	}()
}

// DebugStatus dumps *controllers.MediaStatusResponse
func (c *controller) DebugStatus(status *controllers.MediaStatusResponse) {
	c.logger.Debug("DebugStatus(): *ctl.MediaStatusResponse:status")
	grok.Value(status)
}
