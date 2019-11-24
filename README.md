# go-shodan
[![Build Status](https://travis-ci.org/ns3777k/go-shodan.svg?branch=master)](https://travis-ci.org/ns3777k/go-shodan)
[![Build status](https://ci.appveyor.com/api/projects/status/wbi5u34k5pokbypj/branch/master?svg=true)](https://ci.appveyor.com/project/ns3777k/go-shodan/branch/master)
[![codecov](https://codecov.io/gh/ns3777k/go-shodan/branch/master/graph/badge.svg)](https://codecov.io/gh/ns3777k/go-shodan)
[![GoDoc](https://godoc.org/github.com/ns3777k/go-shodan/shodan?status.svg)](https://godoc.org/github.com/ns3777k/go-shodan/shodan)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ns3777k/go-shodan)](https://goreportcard.com/report/github.com/ns3777k/go-shodan)

To start working with Shodan you have to get your token first. You can do this at [https://www.shodan.io](https://www.shodan.io).

### Usage

The import path depends on whether you use go modules:

```go
import "github.com/ns3777k/go-shodan/v4/shodan"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
import "github.com/ns3777k/go-shodan/shodan" // with go modules disabled
```

Simple example of resolving hostnames:

```go
package main

import (
	"log"
	"context"

	"github.com/ns3777k/go-shodan/v4/shodan" // go modules required
)

func main() {
	client := shodan.NewEnvClient(nil)
	dns, err := client.GetDNSResolve(context.Background(), []string{"google.com", "ya.ru"})

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
	"context"

	"github.com/ns3777k/go-shodan/v4/shodan" // go modules required
)

func main() {
	client := shodan.NewEnvClient(nil)
	ch := make(chan *shodan.HostData)
	err := client.GetBannersByASN(context.Background(), []string{"3303", "32475"}, ch)
	if err != nil {
		panic(err)
	}

	for {
		banner, ok := <-ch

		if !ok {
			log.Println("channel was closed")
			break
		}

		log.Println(banner.Product)
	}
}
```

### Tips and tricks

Every method accepts context in the first argument so you can easily cancel any request.

You can also use `SetDebug(true)` to see the curl version of your requests.

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
- [x] /shodan/alert/triggers
- [x] /shodan/alert/{id}/trigger/{trigger}
- [x] /shodan/alert/{id}/trigger/{trigger}/ignore/{service}

#### Directory Methods
- [x] /shodan/query
- [x] /shodan/query/search
- [x] /shodan/query/tags

#### Account Methods
- [x] /account/profile

#### DNS Methods
- [x] /dns/resolve
- [x] /dns/reverse
- [x] /dns/domain/{domain}

#### Bulk Data
- [x] /shodan/data
- [x] /shodan/data/{dataset}

#### Manage Organization
- [x] /org
- [x] /org/member/{user}

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
