package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type downloader struct {
	url      string
	size     int
	attempts int
	client   *http.Client
}

func newDownloader() *downloader {
	return &downloader{
		url:      os.Getenv("SERVER"),
		size:     100 * 1024 * 1024,
		attempts: 3,
		client:   http.DefaultClient,
	}
}

func (this *downloader) printableResult() string {
	return fmt.Sprintf("Avg. download speed: %.3fMb/s", this.run())
}

func (this *downloader) run() float32 {
	results := make(chan float32, this.attempts)

	for i := 0; i < this.attempts; i++ {
		go func() {
			url := fmt.Sprintf("%s/download?size=%d", this.url, this.size)

			req, _ := http.NewRequest("GET", url, nil)

			// Start timing.
			t := time.Now().UnixNano()

			// Run request
			resp, err := this.client.Do(req)

			if err != nil {
				log.Println(err)
				results <- 0
			}

			// Stop timing (errors might still have happened".
			t = time.Now().UnixNano() - t

			if resp.StatusCode != 200 {
				log.Println(resp.Status)
				results <- 0
			}

			// Read body so that trailer headers become accessible
			resp.Body.Read(nil)

			// Subtract process time
			h := resp.Trailer.Get("X-Duration")
			processTime, err := strconv.ParseInt(h, 0, 0)

			if err != nil {
				log.Printf("Bad download header value: %s\n", h)
				results <- 0
			}

			t -= processTime

			// Push result
			results <- float32(this.size) / float32(t)
		}()

		time.Sleep(1 * time.Second)
	}

	var r float32
	var sum float32 = 0
	successfulAttempts := this.attempts

	for i := 0; i < this.attempts; i++ {
		r = <-results

		if r == 0 {
			successfulAttempts -= 1
			continue
		}

		sum += r
	}

	if successfulAttempts == 0 {
		return 0
	} else {
		return sum / float32(successfulAttempts)
	}
}
