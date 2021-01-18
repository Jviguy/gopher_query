# gopher_query
A simple yet decently fast minecraft query library.
## Installation
```bash
go get github.com/jviguy/gopher_query
```
## Usage
```go
import "github.com/jviguy/gopher_query"
import "log"
import "fmt"

func main() {
  c := gopher_query.NewClient()
  data, err := c.LongQuery("velvetpractice.live:19132")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(data.Players)
}```
