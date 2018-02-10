package main

import (
	gglh "github.com/hiromaily/go-google-home"
	lg "github.com/hiromaily/golibs/log"
	"log"
	"time"
)

func init() {
	//log
	lg.InitializeLog(lg.DebugStatus, lg.LogOff, log.Lshortfile,
		"[Google-Home]", "")
}

func main() {
	// 1.discover Google Home
	gh := gglh.DiscoverService()
	gh.NewClient()

	// 2.speak something
	defer gh.Close()
	err := gh.Speak("This is first speaking test.", "en")
	if err != nil {
		lg.Errorf("gh.Speak() error:%v", err)
	}
	time.Sleep(5 * time.Second)
}
