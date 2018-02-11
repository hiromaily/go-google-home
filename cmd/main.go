package main

import (
	"flag"
	gglh "github.com/hiromaily/go-google-home"
	lg "github.com/hiromaily/golibs/log"
	"log"
)

var (
	message  = flag.String("msg", "", "Message to Google Home")
	lang     = flag.String("lang", "en", "Language to speak")
	server   = flag.Bool("server", false, "Run by server mode")
	logLevel = flag.Int("log", 1, "Run by debug mode")
)

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

	} else {
		// 4.speak something
		err := gh.Speak(*message, *lang)
		if err != nil {
			lg.Errorf("gh.Speak() error:%v", err)
			return
		}
	}
	//time.Sleep(5 * time.Second)
}
