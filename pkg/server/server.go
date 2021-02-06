package server

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"os"
//	"os/signal"
//	"time"
//
//	lg "github.com/hiromaily/golibs/log"
//)
//
//// StartServer is to start server mode
//func (g *GoogleHome) StartServer(port int, lang string) {
//	stopCh := make(chan os.Signal, 1)
//	signal.Notify(stopCh, os.Interrupt)
//	defer signal.Stop(stopCh)
//
//	//server object
//	ghs := GHServer{
//		lang:       lang,
//		GoogleHome: g,
//	}
//
//	http.HandleFunc("/speak", ghs.speakHandler())
//	//http.Handle("/ssml/", http.StripPrefix("/ssml/", http.FileServer(http.Dir("./ssml"))))
//
//	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: http.DefaultServeMux}
//	lg.Infof("Server start with port %d ...", port)
//
//	go func() {
//		// service connections
//		if err := srv.ListenAndServe(); err != nil {
//			lg.Infof("listen: %s\n", err)
//			return
//		}
//	}()
//
//	<-stopCh // wait for SIGINT
//
//	lg.Info("Shutting down server...")
//	// shut down gracefully, but wait no longer than 5 seconds before halting
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	srv.Shutdown(ctx)
//	lg.Info("Server gracefully stopped")
//}
//
////-----------------------------------------------------------------------------
//// GHServer
////-----------------------------------------------------------------------------
//
//// GHServer is google home server
//type GHServer struct {
//	lang string
//	*GoogleHome
//}
//
//// Speak is test object
//type Speak struct {
//	Text string `json:"text"`
//}
//
//func (g *GHServer) speakHandler() func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		//check post or not
//		if r.Method != "POST" {
//			http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
//			return
//		}
//		//check parameter in json
//		text, err := parseJSON(w, r)
//		if err != nil {
//			return
//		}
//
//		err = g.speak(w, text)
//		if err != nil {
//			return
//		}
//
//		lg.Infof("said: %s", text)
//
//		//response correctly
//		w.WriteHeader(http.StatusOK)
//		fmt.Fprint(w, text)
//	}
//}
//
//func (g *GHServer) speak(w http.ResponseWriter, text string) error {
//	if text != "" {
//		err := g.Speak(text, g.lang)
//		if err != nil {
//			http.Error(w, "it couldn't speak", http.StatusInternalServerError)
//			lg.Errorf("gh.Speak() error:%v", err)
//			return err
//		}
//	} else {
//		http.Error(w, "parameter is invalid", http.StatusBadRequest)
//		lg.Error("gh.Speak() error: text is blank")
//		return fmt.Errorf("parameter is invalid")
//	}
//	return nil
//}
//
//func parseJSON(w http.ResponseWriter, r *http.Request) (string, error) {
//	b, err := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		lg.Errorf("gh.parseJson() error:%v", err)
//		return "", err
//	}
//
//	var speak Speak
//	err = json.Unmarshal(b, &speak)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		lg.Errorf("gh.parseJson() error:%v", err)
//		return "", err
//	}
//
//	return speak.Text, err
//}
