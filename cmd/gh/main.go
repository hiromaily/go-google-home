package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-google-home/pkg/commands"
	"github.com/hiromaily/go-google-home/pkg/config"
	"github.com/hiromaily/go-google-home/pkg/files"
)

var (
	tomlPath = flag.String("toml", "", "TOML file path")
	addr     = flag.String("addr", "", "device address")
	lang     = flag.String("lang", "", "language to speak")
	volume   = flag.String("vol", "", "volume: 0.0-1.0")

	//message = flag.String("msg", "", "Message to Google Home")
	//music   = flag.String("music", "", "URL of Music file")
	//server  = flag.Bool("server", false, "Run by server mode")
	// serverPort = flag.Int("port", 8080, "Server port")
)

//var usage = `Usage: %s [options...]
//Options:
//  -msg    Message to Google Home.
//  -music  URL of Music file.      [e.g.] http://music.xxx/music.mp3
//  -lang   Language to speak.      [e.g.] en, de, nl, fr, ja ...
//  -vol    Volume: 0.0-1.0         [e.g.] 0.3
//  -server Run by server mode.     [e.g.] $ gh -server
//  -port   Server port.
//
//See Makefile for examples.
//`

func parseFlag() {
	//flag.Usage = func() {
	//	fmt.Fprintf(os.Stderr, usage, os.Args[0])
	//}
	//flag.Parse()
	flag.Parse()
}

//func checkVersion() {
//	if *isVersion {
//		fmt.Printf("%s %s\n", "gh", version)
//		os.Exit(0)
//	}
//}

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

	conf := getConfig()
	// overwrite config
	if *addr != "" {
		conf.Device.Address = *addr
	}
	if *lang != "" {
		conf.Device.Lang = *lang
	}

	regi := NewRegistry(conf)
	devicer, err := regi.NewDevicer().Start()
	if err != nil {
		log.Fatalf("fail to connect Google Home: %v", err)
		return
	}

	// wait events
	chFinishNotifier := make(chan struct{})
	defer func() {
		devicer.Close()
		close(chFinishNotifier)
	}()

	// register sub commands
	commands.Register(regi.NewLogger(), devicer, chFinishNotifier)
	os.Exit(int(subcommands.Execute(context.Background())))

	// volume TODO:Fix DATA RACE
	//if *volume != "" {
	//	devicer.Controller().SetVolume(*volume)
	//}

	//switch {
	//case *server:
	//	// server mode
	//	// gh.StartServer(*serverPort, *lang)
	//	return
	//case *music != "":
	//	log.Printf("play: %s", *music)
	//	gh.Controller().RunEventReceiver(finishNotification)
	//
	//	// play music
	//	err = gh.Controller().Play(*music)
	//}

	//monitorStatus()
	<-chFinishNotifier
}

//func monitorStatus() {
//	//TODO: It causes DATA RACE
//	//m := new(sync.Mutex)
//	//go func() {
//	//	status, err := gh.GetStatus()
//	//	if err != nil {
//	//		lg.Errorf("gh.GetStatus() error:%v", err)
//	//		return
//	//	} else {
//	//		m.Lock()
//	//		gh.DebugStatus(status)
//	//		m.Unlock()
//	//	}
//	//}()
//}
