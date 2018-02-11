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
	"time"
)

var (
	message    = flag.String("msg", "", "Message to Google Home")
	lang       = flag.String("lang", "en", "Language to speak")
	server     = flag.Bool("server", false, "Run by server mode")
	serverPort = flag.Int("port", 8080, "Server por")
	logLevel   = flag.Int("log", 1, "Run by debug mode")
)

type GHServer struct {
	*gglh.GoogleHome
}

type Speak struct {
	Text string `json:"text"`
}

func init() {
	flag.Parse()

	//log
	lg.InitializeLog(uint8(*logLevel), lg.LogOff, log.Lshortfile,
		"[Google-Home]", "")
}

func main() {
	if !*server && *message == "" {
		lg.Error("Please type in msg option.")
		return
	}

	// 1.discover Google Home
	gh := gglh.DiscoverService()
	if gh.Error != nil {
		lg.Errorf("gglh.DiscoverService() error:%v", gh.Error)
		return
	}
	// if you use specific address
	//gh := gglh.New("192.168.178.164", 8009)

	// 2.create client
	gh.NewClient()
	defer gh.Close()

	// 3.server mode
	if *server {
		listen(gh)
	} else {
		// 4.speak something
		err := gh.Speak(*message, *lang)
		if err != nil {
			lg.Errorf("gh.Speak() error:%v", err)
			return
		}
	}
	time.Sleep(1 * time.Second)
}

func listen(gh *gglh.GoogleHome) {
	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt)

	//server object
	ghs := GHServer{}
	ghs.GoogleHome = gh

	http.HandleFunc("/speak", ghs.handler())
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

func (g *GHServer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//check post or not
		if r.Method != "POST" {
			http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
			return
		}
		//check parameter
		r.ParseForm()
		var said string
		var err error
		if v, ok := r.Form["text"]; ok {
			said, err = g.speak(w, v[0])
		} else {
			//json
			said, err = g.parseJson(w, r)
		}
		if err != nil {
			return
		}
		lg.Infof("said: %s", said)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, said)
	}
}

func (g *GHServer) speak(w http.ResponseWriter, text string) (string, error) {
	if text != "" {
		err := g.Speak(text, *lang)
		if err != nil {
			http.Error(w, "It couldn't speak.", http.StatusInternalServerError)
			lg.Errorf("gh.Speak() error:%v", err)
			return "", err
		}
	} else {
		http.Error(w, "Parameter is invalid.", http.StatusBadRequest)
		lg.Error("gh.Speak() error: text is blank")
		return "", fmt.Errorf("Parameter is invalid.")
	}
	return text, nil
}

func (g *GHServer) parseJson(w http.ResponseWriter, r *http.Request) (string, error) {
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

	return g.speak(w, speak.Text)
}
