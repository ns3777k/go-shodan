# go-shodan
To start working with Shodan you have to get your token first. You can do this at [http://www.shodanhq.com](http://www.shodanhq.com).
#### Please note that it is unstable for now

### Installation
Download the package:

```bash
go get "github.com/ns3777k/go-shodan/shodan"
```

Start using it:
```go
package main

import (
    "log"
    
    "github.com/ns3777k/go-shodan/shodan"
)

func main() {
    client := shodan.NewClient("MY_TOKEN")
    dns, err := client.GetDnsResolve([]string{"google.com", "ya.ru"})
    
    if err != nil {
        log.Panic(err)
    } else {
        log.Println(dns["google.com"])
    }
}
```
Sample output:
```bash
$ go run c.go
2015/09/05 18:50:52 173.194.115.35
```

### Roadmap
1. Streaming API
2. Error handling
3. Testing

### Links
* [goDoc](http://godoc.org/github.com/ns3777k/go-shodan/shodan)
* [Shodan.io](http://shodan.io) & [ShodanHQ](http://www.shodanhq.com)
* [API Documentation](https://developer.shodan.io/api)
