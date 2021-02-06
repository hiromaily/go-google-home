package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-google-home/pkg/config"
	"github.com/hiromaily/go-google-home/pkg/files"
)

var (
	tomlPath = flag.String("toml", "", "TOML file path")
	message  = flag.String("msg", "", "Message to Google Home")
	music    = flag.String("music", "", "URL of Music file")
	// address    = flag.String("addr", os.Getenv("GOOGLE_HOME_IP"), "Address of Google Home (e.g. 192.168.x.x:8009)")
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

func parseFlag() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}
	flag.Parse()
}

//func checkVersion() {
//	if *isVersion {
//		fmt.Printf("%s %s\n", "book-teacher", version)
//		os.Exit(0)
//	}
//}

func validateArguments() {
	// this pattern is not allowed
	if !*server && *message == "" && *music == "" {
		flag.Usage()
		os.Exit(1)
		return
	}
}

func getConfig() *config.Root {
	configPath := files.GetConfigPath(*tomlPath)
	if configPath == "" {
		log.Fatal(errors.New("config file is not found"))
	}
	log.Println("config file: ", configPath)
	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	return conf
}

func main() {
	parseFlag()
	validateArguments()

	conf := getConfig()
	regi := NewRegistry(conf)
	devicer := regi.NewDevicer()
	gh, err := devicer.Start(*address)
	if err != nil {
		log.Fatalf("fail to connect Google Home: %v", err)
		return
	}
	defer gh.Controller().Close()

	// volume TODO:Fix DATA RACE
	if *volume != "" {
		gh.Controller().SetVolume(*volume)
	}

	// wait events
	finishNotification := make(chan bool)

	switch {
	case *server:
		// server mode
		// gh.StartServer(*serverPort, *lang)
		return
	case *message != "":
		log.Printf("speak: %s", *message)
		gh.Controller().RunEventReceiver(finishNotification)

		// speak something
		err = gh.Controller().Speak(*message, *lang)
	case *music != "":
		log.Printf("play: %s", *music)
		gh.Controller().RunEventReceiver(finishNotification)

		// play music
		err = gh.Controller().Play(*music)
	default:
	}
	if err != nil {
		log.Printf("fail to speak/play: %v", err)
		close(finishNotification)
		gh.Controller().CloseEvent()
		return
	}

	monitorStatus()

	<-finishNotification
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
