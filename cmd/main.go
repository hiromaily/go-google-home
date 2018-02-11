package main

import (
	"flag"
	gglh "github.com/hiromaily/go-google-home"
	lg "github.com/hiromaily/golibs/log"
	"log"
)

var (
	message = flag.String("msg", "Please type message as option.", "Message to Google Home")
	lang    = flag.String("lang", "en", "Language to speak")
	server  = flag.Bool("server", false, "Run by server mode")
)

func init() {
	flag.Parse()

	//log
	lg.InitializeLog(lg.DebugStatus, lg.LogOff, log.Lshortfile,
		"[Google-Home]", "")
}

func main() {
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

	// 3.speak something
	defer gh.Close()
	err := gh.Speak(*message, *lang)
	if err != nil {
		lg.Errorf("gh.Speak() error:%v", err)
	}
	//time.Sleep(5 * time.Second)
}
