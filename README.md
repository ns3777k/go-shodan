# go-shodan
[![Build Status](https://travis-ci.org/ns3777k/go-shodan.svg?branch=master)](https://travis-ci.org/ns3777k/go-shodan)
[![Coverage Status](https://coveralls.io/repos/ns3777k/go-shodan/badge.svg?branch=master&service=github)](https://coveralls.io/github/ns3777k/go-shodan?branch=master)
[![GoDoc](https://godoc.org/github.com/ns3777k/go-shodan/shodan?status.svg)](https://godoc.org/github.com/ns3777k/go-shodan/shodan)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ns3777k/go-shodan)](https://goreportcard.com/report/github.com/ns3777k/go-shodan)

To start working with Shodan you have to get your token first. You can do this at [https://www.shodan.io](https://www.shodan.io).

### Installation
Download the package:

```bash
go get "github.com/ns3777k/go-shodan/shodan"
```

That's it. You're ready to roll :-)

### Usage

Simple example of resolving hostnames:

```go
package main

import (
    "log"
    
    "github.com/ns3777k/go-shodan/shodan"
)

func main() {
    client := shodan.NewClient(nil, "MY_TOKEN")
    dns, err := client.GetDNSResolve([]string{"google.com", "ya.ru"})
    
    if err != nil {
        log.Panic(err)
    } else {
        log.Println(*dns["google.com"])
    }
}
```
Output for above:
```bash
2015/09/05 18:50:52 173.194.115.35
```

Streaming example:

```go
package main

import (
    "log"
    
    "github.com/ns3777k/go-shodan/shodan"
)

func main() {
    client := shodan.NewClient(nil, "MY_TOKEN")

    go func() {
        for {
            banner, ok := <- client.StreamChan
            if !ok {
                log.Fatalln("channel got closed")
            }

            // Do something here with banner
        }
    }()

    go client.GetBanners()

    for {
        time.Sleep(time.Second * 10)
    }
}
```

### Implemented API

#### Search Methods
- [x] /shodan/host/{ip}
- [x] /shodan/host/count
- [x] /shodan/host/search
- [x] /shodan/host/search/tokens
- [x] /shodan/ports

#### On-Demand Scanning
- [x] /shodan/protocols
- [x] /shodan/scan
- [x] /shodan/scan/internet
- [x] /shodan/scan/{id}

#### Network Alerts
- [ ] /shodan/alert
- [ ] /shodan/alert/{id}/info
- [ ] /shodan/alert/{id}
- [ ] /shodan/alert/info

#### Directory Methods
- [x] /shodan/query
- [x] /shodan/query/search
- [x] /shodan/query/tags

#### Account Methods
- [x] /account/profile

#### DNS Methods
- [x] /dns/resolve
- [x] /dns/reverse

#### Utility Methods
- [x] /tools/httpheaders
- [x] /tools/myip

#### API Status Methods
- [x] /api-info

#### Experimental Methods
- [x] /labs/honeyscore/{ip}

If a method is absent or something doesn't work properly don't hesitate to create an issue.

### Links
* [Shodan.io](http://shodan.io) & [ShodanHQ](http://www.shodanhq.com)
* [API Documentation](https://developer.shodan.io/api)
