# go-google-home
It makes Google Home speak something and is inspired by 
[google-home-notifier](https://github.com/noelportugal/google-home-notifier).  
There are 2 modes.
- command line mode with text parameter.
- server mode hadling post message.


## Sample code
It is in `cmd` directory with Makefile.

```
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
```


#### About options in cmd/main.go
| Options        |                                           | Type   | Example                 |
| -------------- | ------------------------------------------ | -------| ---------------------- |
| msg            | Message to Google Home                     | string | "Hello world!"         |
| addr           | IP address + Port for specific Google Home | string | "xxx.xxx.xxx.xxx:8009" |
| lang           | Language to speak                          | string | en                     |
| server         | Run by server mode                         | bool   | none                   |
| port           | Web Server port                            | int    | 8080                   |
| log            | Log level, `1` displays even debug message | int    | 1                      |

#### Execution example
```
# Build
$ go build -i -race -v -o ${GOPATH}/bin/gh ./cmd/


# Sample for saying something in English.
$ gh -msg "Thank you." -lang en

# Sample for saying something in Japanese.
$ gh -msg "ありがとうございます" -lang ja

# Sample for saying something in French.
$ gh -msg "Merci." -lang fr

# Sample for saying something in German.
$ gh -msg "Danke." -lang de


# Sample for saying something in English with `debug` log.
$ gh -msg "This displays debug log." -log 1


# Sample for using specific IP address of Google Home.
$ gh  -msg "It reaches to specific IP address." -addr "10.0.0.1:8009"


# Sample for server mode.
$ gh -server

# Sample to post message to server by HTTPie
$ http POST http://localhost:8080/speak text="It's sunny day today."
```

## How to access to local server from outside
Use [Ngrok](https://github.com/inconshreveable/ngrok)

#### Install on Mac
```
$ brew cask install ngrok
```
```
# If you use 8080 port for that local server.
$ ngrok http 8080
```
