package main

import (
	"flag"
	"fmt"
	"os"

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

	gh := createClient()
	if gh.Error != nil {
		lg.Errorf("fail to connect Google Home: %v", gh.Error)
		return
	}
	defer gh.Close()

	//volume TODO:Fix DATA RACE
	if *volume != "" {
		gh.SetVolume(*volume)
	}

	// wait events
	finishNotification := make(chan bool)
	var err error

	switch {
	case *server:
		// server mode
		gh.StartServer(*serverPort, *lang)
		return
	case *message != "":
		lg.Infof("speak: %s", *message)
		gh.RunEventReceiver(finishNotification)

		// speak something
		err = gh.Speak(*message, *lang)
	case *music != "":
		lg.Infof("play: %s", *music)
		gh.RunEventReceiver(finishNotification)

		// play music
		err = gh.Play(*music)
	default:
	}
	if err != nil {
		lg.Errorf("fail to speak/play: %v", err)
		close(finishNotification)
		close(gh.Client.Events)
		return
	}

	monitorStatus()

	<-finishNotification
}

func createClient() *gglh.GoogleHome {
	var gh *gglh.GoogleHome

	//TODO: is it better to environment variable if existing?
	//os.Getenv("GOOGLE_HOME_IP")
	if *address != "" {
		// create object from address
		lg.Infof("from address: %s", *address)
		gh = gglh.NewGoogleHome().WithAddressString(*address).WithClient()
	} else {
		// discover Google Home
		lg.Info("discover google home address")
		gh = gglh.DiscoverService().WithClient()
	}
	return gh
}

func monitorStatus() {
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
}
