package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	gglh "github.com/hiromaily/go-google-home"
	lg "github.com/hiromaily/golibs/log"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	//"sync"
	"time"
)

var (
	message    = flag.String("msg", "", "Message to Google Home")
	address    = flag.String("addr", os.Getenv("GOOGLE_HOME_IP"), "Address of Google Home (e.g. 192.168.x.x:8009)")
	lang       = flag.String("lang", "en", "Language to speak")
	server     = flag.Bool("server", false, "Run by server mode")
	serverPort = flag.Int("port", 8080, "Server port")
	logLevel   = flag.Int("log", 2, "log level. `1`:debug mode")
)

var usage = `Usage: %s [options...]
Options:
  -msg    Message to Google Home.
  -addr   Address of Google Home. [e.g.] 192.168.x.x:8009
  -lang   Language to speak.      [e.g.] en, de, nl, fr, ja ...
  -server Run by server mode.     [e.g.] $ gh -server
  -port   Server port.
  -log    Log level.              [e.g.] 1: debug log

See Makefile for examples.
`

type GHServer struct {
	*gglh.GoogleHome
}

type Speak struct {
	Text string `json:"text"`
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}

	flag.Parse()

	//log
	var logFmt = log.Lshortfile
	if *logLevel != 1 {
		logFmt = 0
	}
	lg.InitializeLog(uint8(*logLevel), lg.LogOff, logFmt,
		"[Google-Home]", "")
}

func main() {
	//m := new(sync.Mutex)

	if !*server && *message == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	var gh *gglh.GoogleHome
	var err error

	if *address != "" {
		// create object from address
		gh, err = setAddress(*address)
		if err != nil {
			lg.Error(err)
			return
		}
	} else {
		// discover Google Home
		gh = gglh.DiscoverService()
		if gh.Error != nil {
			lg.Errorf("gglh.DiscoverService() error:%v", gh.Error)
			return
		}
	}

	// create client
	gh.NewClient()
	defer gh.Close()

	// server mode
	if *server {
		listen(gh)
		return
	}

	// wait events
	finishNotification := make(chan bool)
	gh.RunEventReceiver(finishNotification)

	// speak something
	err = gh.Speak(*message, *lang)
	if err != nil {
		lg.Errorf("gh.Speak() error:%v", err)
		return
	}

	//TODO: monitor status
	//TODO: it causes DATA RACE
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

func setAddress(address string) (*gglh.GoogleHome, error) {
	// 1. use address if it exists.
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		return nil, fmt.Errorf("addr argument is invalid. It should be :%s", "xxx.xxx.xxx.xxx:8009")
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, fmt.Errorf("addr argument is invalid. It should be :%s", "xxx.xxx.xxx.xxx:8009")
	}
	return gglh.New(addr[0], port), nil
}

func listen(gh *gglh.GoogleHome) {
	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt)

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
		text, err := parseJson(w, r)
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

//func (g *GHServer) ssmlHandler() func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		//url
//		url := "http://192.168.178.157:8080/ssml/sample.ssml"
//		err := g.SpeakBySSML(url)
//		if err != nil {
//			http.Error(w, "It couldn't speak.", http.StatusInternalServerError)
//			lg.Errorf("gh.Speak() error:%v", err)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		fmt.Fprint(w, 200)
//	}
//}

func (g *GHServer) speak(w http.ResponseWriter, text string) error {
	if text != "" {
		err := g.Speak(text, *lang)
		if err != nil {
			http.Error(w, "It couldn't speak.", http.StatusInternalServerError)
			lg.Errorf("gh.Speak() error:%v", err)
			return err
		}
	} else {
		http.Error(w, "Parameter is invalid.", http.StatusBadRequest)
		lg.Error("gh.Speak() error: text is blank")
		return fmt.Errorf("Parameter is invalid.")
	}
	return nil
}

func parseJson(w http.ResponseWriter, r *http.Request) (string, error) {
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
