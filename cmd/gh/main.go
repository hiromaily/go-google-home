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
)

//var usage = `Usage: %s [options...]
//Options:
//  -toml  TOML file path
//  -addr  device address
//  -lang  language to speak e.g.) en, de, nl, fr, ja ...
//  -vol   volume: 0.0-1.0   e.g.) 0.5
//`

func parseFlag() {
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

	// volume TODO:Fix DATA RACE
	//if *volume != "" {
	//	devicer.Controller().SetVolume(*volume)
	//}

	// wait events
	chFinishNotifier := make(chan struct{})
	var commandResult int
	defer func() {
		devicer.Close()
		close(chFinishNotifier)
		os.Exit(commandResult)
	}()

	// register sub commands
	commands.Register(regi.NewLogger(), devicer, chFinishNotifier)
	// execute sub command
	commandResult = int(subcommands.Execute(context.Background()))

	<-chFinishNotifier
}
