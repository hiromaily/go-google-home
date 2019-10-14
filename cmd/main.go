package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	gglh "github.com/hiromaily/go-google-home/pkg/googlehome"
	lg "github.com/hiromaily/golibs/log"
)

var (
	message = flag.String("msg", "", "Message to Google Home")
	music   = flag.String("music", "", "URL of Music file")
	//address    = flag.String("addr", os.Getenv("GOOGLE_HOME_IP"), "Address of Google Home (e.g. 192.168.x.x:8009)")
	address    = flag.String("addr", "", "Address of Google Home (e.g. 192.168.x.x:8009)")
	lang       = flag.String("lang", "en", "Language to speak")
	volume     = flag.String("vol", "", "Volume: 0.0-1.0")
	server     = flag.Bool("server", false, "Run by server mode")
	serverPort = flag.Int("port", 8080, "Server port")
)

var usage = `Usage: %s [options...]
Options:
  -msg    Message to Google Home.
  -music  URL of Music file.      [e.g.] http://music.xxx/music.mp3
  -addr   Address of Google Home. [e.g.] 192.168.x.x:8009
  -lang   Language to speak.      [e.g.] en, de, nl, fr, ja ...
  -vol    Volume: 0.0-1.0         [e.g.] 0.3 
  -server Run by server mode.     [e.g.] $ gh -server
  -port   Server port.
  -log    Log level.              [e.g.] 1: debug log

See Makefile for examples.
`

// GHServer is google home server
type GHServer struct {
	*gglh.GoogleHome
}

// Speak is test object
type Speak struct {
	Text string `json:"text"`
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}
	flag.Parse()
}

func validateArguments() {
	// this pattern is not allowed
	if !*server && *message == "" && *music == "" {
		flag.Usage()
		os.Exit(1)
		return
	}
}

func main() {
	//validate
	validateArguments()

	lg.InitializeLog(lg.DebugStatus, lg.TimeShortFile, "[Google-Home]", "", "hiromaily")

	var gh *gglh.GoogleHome
	var err error

	//TODO: is it better to environment variable if existing
	//os.Getenv("GOOGLE_HOME_IP")
	if *address != "" {
		// create object from address
		lg.Infof("from address: %s", *address)
		gh, err = gglh.NewGoogleHome().WithAddressString(*address)
		if err != nil {
			lg.Error(err)
			return
		}
	} else {
		// discover Google Home
		lg.Info("discover google home address")
		gh = gglh.DiscoverService()
		if gh.Error != nil {
			lg.Errorf("gglh.DiscoverService() error:%v", gh.Error)
			return
		}
	}

	// create client
	gh.NewClient()
	defer gh.Close()

	//volume
	//TODO:Fix DATA RACE
	if *volume != "" {
		gh.SetVolume(*volume)
	}

	// server mode
	if *server {
		listen(gh)
		return
	}

	// wait events
	finishNotification := make(chan bool)
	gh.RunEventReceiver(finishNotification)

	if *message != "" {
		// speak something
		err = gh.Speak(*message, *lang)
		if err != nil {
			lg.Errorf("gh.Speak() error:%v", err)
			close(finishNotification)
			close(gh.Client.Events)
			return
		}
	} else if *music != "" {
		// play music
		err = gh.Play(*music)
		if err != nil {
			lg.Errorf("gh.Play() error:%v", err)
			close(finishNotification)
			close(gh.Client.Events)
			return
		}
	} else {
		//this part should not be passed.
		close(finishNotification)
		return
	}

	//TODO: monitor status
	//TODO: It causes DATA RACE
	//m := new(sync.Mutex)
	//go func() {
	//	status, err := gh.GetStatus()
	//	if err != nil {
	//		lg.Errorf("gh.GetStatus() error:%v", err)
	//		return
	//	} else {
	//		m.Lock()
	//		gh.DebugStatus(status)
	//		m.Unlock()
	//	}
	//}()

	<-finishNotification
}

func listen(gh *gglh.GoogleHome) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)
	defer signal.Stop(stopCh)

	//server object
	ghs := GHServer{}
	ghs.GoogleHome = gh

	http.HandleFunc("/speak", ghs.speakHandler())
	//http.Handle("/ssml/", http.StripPrefix("/ssml/", http.FileServer(http.Dir("./ssml"))))

	srv := &http.Server{Addr: fmt.Sprintf(":%d", *serverPort), Handler: http.DefaultServeMux}
	lg.Infof("Server start with port %d ...", *serverPort)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			lg.Infof("listen: %s\n", err)
			return
		}
	}()

	<-stopCh // wait for SIGINT

	lg.Info("Shutting down server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	lg.Info("Server gracefully stopped")
}

func (g *GHServer) speakHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//check post or not
		if r.Method != "POST" {
			http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
			return
		}
		//check parameter in json
		text, err := parseJSON(w, r)
		if err != nil {
			return
		}

		err = g.speak(w, text)
		if err != nil {
			return
		}

		lg.Infof("said: %s", text)

		//response correctly
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, text)
	}
}

func (g *GHServer) speak(w http.ResponseWriter, text string) error {
	if text != "" {
		err := g.Speak(text, *lang)
		if err != nil {
			http.Error(w, "it couldn't speak", http.StatusInternalServerError)
			lg.Errorf("gh.Speak() error:%v", err)
			return err
		}
	} else {
		http.Error(w, "parameter is invalid", http.StatusBadRequest)
		lg.Error("gh.Speak() error: text is blank")
		return fmt.Errorf("parameter is invalid")
	}
	return nil
}

func parseJSON(w http.ResponseWriter, r *http.Request) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		lg.Errorf("gh.parseJson() error:%v", err)
		return "", err
	}

	var speak Speak
	err = json.Unmarshal(b, &speak)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		lg.Errorf("gh.parseJson() error:%v", err)
		return "", err
	}

	return speak.Text, err
}
