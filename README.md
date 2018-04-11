# go-shodan
[![Build Status](https://travis-ci.org/ns3777k/go-shodan.svg?branch=master)](https://travis-ci.org/ns3777k/go-shodan)
[![Build status](https://ci.appveyor.com/api/projects/status/wbi5u34k5pokbypj/branch/master?svg=true)](https://ci.appveyor.com/project/ns3777k/go-shodan/branch/master)
[![codecov](https://codecov.io/gh/ns3777k/go-shodan/branch/master/graph/badge.svg)](https://codecov.io/gh/ns3777k/go-shodan)
[![GoDoc](https://godoc.org/github.com/ns3777k/go-shodan/shodan?status.svg)](https://godoc.org/github.com/ns3777k/go-shodan/shodan)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ns3777k/go-shodan)](https://goreportcard.com/report/github.com/ns3777k/go-shodan)
[![codebeat badge](https://codebeat.co/badges/12e593ad-ca40-41e8-9b84-61316947d2eb)](https://codebeat.co/projects/github-com-ns3777k-go-shodan-master)

To start working with Shodan you have to get your token first. You can do this at [https://www.shodan.io](https://www.shodan.io).

### Installation
Download the package:

```bash
go get "gopkg.in/ns3777k/go-shodan.v2"
```

To use the old version:

```bash
go get "gopkg.in/ns3777k/go-shodan.v1"
```

That's it. You're ready to roll :-)

`master`-branch is kind of stable but might have some breaking changes. See the changelog for details.

### Usage

Simple example of resolving hostnames:

```go
package main

import (
    "log"

    "gopkg.in/ns3777k/go-shodan.v2/shodan"
)

func main() {
    client := shodan.NewEnvClient(nil)
    dns, err := client.GetDNSResolve([]string{"google.com", "ya.ru"})

    if err != nil {
        log.Panic(err)
    } else {
        log.Println(dns["google.com"])
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

    "gopkg.in/ns3777k/go-shodan.v2/shodan"
)

func main() {
    client := shodan.NewEnvClient(nil)

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

### Implemented REST API

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
- [x] /shodan/alert
- [x] /shodan/alert/{id}/info
- [x] /shodan/alert/{id}
- [x] /shodan/alert/info

#### Directory Methods
- [x] /shodan/query
- [x] /shodan/query/search
- [x] /shodan/query/tags

#### Account Methods
- [x] /account/profile

#### DNS Methods
- [x] /dns/resolve
- [x] /dns/reverse

#### Bulk Data
- [ ] /shodan/data
- [ ] /shodan/data/{dataset}

#### Utility Methods
- [x] /tools/httpheaders
- [x] /tools/myip

#### API Status Methods
- [x] /api-info

#### Experimental Methods
- [x] /labs/honeyscore/{ip}

### Implemented Streaming API

#### Data Streams
- [x] /shodan/banners
- [x] /shodan/asn/{asn}
- [x] /shodan/countries/{countries}
- [x] /shodan/ports/{ports}

#### Network Alerts
- [x] /shodan/alert
- [x] /shodan/alert/{id}

If a method is absent or something doesn't work properly don't hesitate to create an issue.

### Links
* [Shodan.io](http://shodan.io)
* [API Documentation](https://developer.shodan.io/api)
