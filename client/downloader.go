package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type result float32

const NO_RESULT result = 0

type downloader struct {
	url     string
	size    int64
	threads int
	client  *http.Client
}

func newDownloader() *downloader {
	return &downloader{
		url:     os.Getenv("SERVER"),
		size:    20 * 1024 * 1024,
		threads: 4,
		client:  http.DefaultClient,
	}
}

func (this *downloader) printableResult() string {
	return fmt.Sprintf("Avg. download speed: %.3fMb/s", this.run())
}

func (this *downloader) run() float32 {
	results := make(chan result, this.threads)

	for i := 0; i < this.threads; i++ {
		go func() {
			url := fmt.Sprintf("%s/download?size=%d", this.url, this.size)

			req, _ := http.NewRequest("GET", url, nil)

			// Start timing.
			t := time.Now().UnixNano()

			// Run request
			resp, err := this.client.Do(req)

			if err != nil {
				log.Println(err)
				results <- NO_RESULT

				return
			}

			if resp.StatusCode != 200 {
				log.Println(resp.Status)
				results <- NO_RESULT

				return
			}

			// Read body so that trailer headers become accessible
			ioutil.ReadAll(resp.Body)

			// Stop timing (errors might still have happened".
			t = time.Now().UnixNano() - t

			// Subtract process time
			h := resp.Trailer.Get("X-Duration")
			processTime, err := strconv.ParseInt(h, 0, 0)

			log.Printf("Process time: %d. Total time: %d", processTime, t)

			if err != nil {
				log.Printf("Bad download header value: %s\n", h)
				results <- NO_RESULT

				return
			}

			t -= processTime

			mbSize := bToMb(this.size)
			sTime := nanoToS(t)

			log.Printf("Downloaded %.3f MB in %.3fs", mbSize, sTime)

			// Push result
			results <- result(mbSize / sTime)
		}()
	}

	var r float32
	var sum float32 = 0

	for i := 0; i < this.threads; i++ {
		r = float32(<-results)

		sum += r
	}

	return sum
}

func bToMb(bytes int64) float32 {
	return float32(bytes) / (1 << 20)
}

func nanoToS(nano int64) float32 {
	return float32(nano) / 1000000000
}
