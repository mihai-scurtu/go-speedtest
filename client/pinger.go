package client

import (
	"net/http"
	"os"
	"time"

	"fmt"

	"golang.org/x/tools/container/intsets"
)

type pinger struct {
	url      string
	client   *http.Client
	attempts int
}

func newPinger() *pinger {
	this := &pinger{
		url:      os.Getenv("SERVER") + "/ping",
		client:   http.DefaultClient,
		attempts: 10,
	}

	return this
}

func (this *pinger) printableResult() string {
	return fmt.Sprintf("Latency: %.3fms", this.run())
}

func (this *pinger) run() float32 {
	var result int64 = int64(intsets.MaxInt)

	ch := make(chan int, this.attempts)

	for i := 0; i < this.attempts; i++ {
		go func() {
			t := time.Now().UnixNano()

			req, _ := http.NewRequest("HEAD", this.url, nil)
			this.client.Do(req)

			t = time.Now().UnixNano() - t

			if t < result {
				result = t
			}

			ch <- i
		}()

		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < this.attempts; i++ {
		<-ch
	}

	return float32(result) / 1000000.0
}
