package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-google-home/pkg/device"
)

//-----------------------------------------------------------------------------
// Server
//-----------------------------------------------------------------------------

// Server interface
type Server interface {
	Start(port int) error
	SpeakHandler() func(http.ResponseWriter, *http.Request)
}

// server object
type server struct {
	logger      *zap.Logger
	devicer     device.Device
	defaultPort int
}

// NewServer returns Server
func NewServer(logger *zap.Logger, devicer device.Device, defaultPort int) Server {
	return &server{
		logger:      logger,
		devicer:     devicer,
		defaultPort: defaultPort,
	}
}

// Start starts server
func (s *server) Start(port int) error {
	if port == 0 {
		port = s.defaultPort
	}

	// device
	_, err := s.devicer.Start()
	if err != nil {
		s.logger.Error("fail to call devicer.Start()", zap.Error(err))
		return err
	}
	defer func() {
		s.devicer.Close()
	}()
	// set callback event
	s.devicer.Controller().RunEventReceiver(nil)

	// signal
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)
	defer signal.Stop(stopCh)

	// handler
	http.HandleFunc("/speak", s.SpeakHandler())
	// http.Handle("/ssml/", http.StripPrefix("/ssml/", http.FileServer(http.Dir("./ssml"))))

	// server
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: http.DefaultServeMux}
	s.logger.Info("Server start", zap.Int("port", port))
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("fail to call ListenAndServe()", zap.Error(err))
			return
		}
	}()
	<-stopCh // wait for SIGINT

	s.logger.Info("shutting down server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info("server gracefully stopped")
	return srv.Shutdown(ctx)
}

// SpeakHandler handles /speak
func (s *server) SpeakHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check post or not
		if r.Method != "POST" {
			http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
			return
		}
		// check parameter in json
		text, err := parseJSON(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.Error("fail to call parseJson()", zap.Error(err))
			return
		}

		err = s.speak(w, text)
		if err != nil {
			return
		}
		s.logger.Info("speak", zap.String("message", text))

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, text)
	}
}

func (s *server) speak(w http.ResponseWriter, text string) error {
	if text == "" {
		msg := "test parameter is invalid"
		http.Error(w, msg, http.StatusBadRequest)
		s.logger.Error(msg)
		return errors.New(msg)
	}
	err := s.devicer.Controller().Speak(text, s.devicer.Lang())
	if err != nil {
		msg := "fail to call Speak()"
		http.Error(w, msg, http.StatusInternalServerError)
		s.logger.Error(msg, zap.Error(err))
		return err
	}
	return nil
}

// Speak object
type Speak struct {
	Text string `json:"text"`
}

func parseJSON(r *http.Request) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", err
	}

	var speak Speak
	err = json.Unmarshal(b, &speak)
	if err != nil {
		return "", err
	}
	return speak.Text, err
}
