# go-google-home

[![Build Status](https://travis-ci.org/hiromaily/go-google-home.svg?branch=master)](https://travis-ci.org/hiromaily/go-google-home)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-google-home)](https://goreportcard.com/report/github.com/hiromaily/go-google-home)
[![codebeat badge](https://codebeat.co/badges/9ddc2e04-f22a-4448-8e7d-ca0c717c76ef)](https://codebeat.co/projects/github-com-hiromaily-go-google-home-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5c83a126d63c402f9a8242295d4a79c4)](https://www.codacy.com/app/hiromaily2/go-google-home?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hiromaily/go-google-home&amp;utm_campaign=Badge_Grade)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-goa/master/LICENSE)


It makes Google Home spoken something and is inspired by [google-home-notifier](https://github.com/noelportugal/google-home-notifier).  
There are 3 modes.
- speak message
- play sound data
- run as server mode then POST message with message can be handled to speak  

### Note
- Google Home device should be in same local network with machine
- Google Home device IP address would be detected automatically. So it's not necessarily to specify device IP address.


## Requirements
- Golang 1.15+
- [direnv](https://github.com/direnv/direnv) for MacOS user
- [Ngrok](https://github.com/inconshreveable/ngrok) if you want to access server from outside

## Installation
### for MacOS user
```
$ brew install hiromaily/tap/go-google-home
 
# config file is installed in /usr/local/etc/google-home/gh.toml
# modify `gh.toml` if settings wanna be changed

# run
$ gh speak -msg "Hi guys, thank you for using. Have fun."
```


## subcommand
| subcommand   |                                                    | example                        |
| ------------ | -------------------------------------------------- | ------------------------------ |
| speak        | speak message on google home                       | gh speak -msg "hello guys"     |
| play         | play sound data on google home                     | gh play -url xxxxx.mp3         |
| server       | run web server to handle http request with message | gh server -port 8888           |


## basic command line option
| options        |                                            | type   | example                      |
| -------------- | ------------------------------------------ | -------| ---------------------------- |
| toml           | TOML file path                             | string | -toml ./configs/default.toml |
| addr           | specify address of Google Homee            | string | -addr xxx.xxx.xxx.xxx:8009   |
| lang           | spoken language, default is english        | string | -lang en                     |
| v              | show version                               | bool   | -v                           |

## environment variable
environment variable `GO_GOOGLE_HOME_CONF` is used as default config path


## example
```
# saying something in English
$ gh speak -msg "Thank you."

# saying something in Japanese
$ gh -lang ja speak -msg "ありがとうございます"

# saying something in Dutch
$ gh -lang nl speak -msg "Dank je" 

# saying something in German
$ gh -lang de speak -msg "Danke."

# saying something in French
$ gh -lang fr speak -msg "Merci."

# playing music
$ gh play -url "https://github.com/hiromaily/go-google-home/raw/master/assets/music/bensound-dubstep.mp3"

# using specific IP address of Google Home.
$ gh -addr "10.0.0.1:8009" -msg "It reaches to specific IP address."
```

## example as server
```
# server mode
$ gh server -port 8888

# then post message to server by HTTPie
$ http POST http://localhost:8080/speak text="It's sunny day today."
```


## How to access to local server from outside easily?
Use [Ngrok](https://github.com/inconshreveable/ngrok)

#### Install on Mac
```
$ brew install --cask ngrok
```

```
# If you use 8080 port for that local server.
$ ngrok http 8080
```
