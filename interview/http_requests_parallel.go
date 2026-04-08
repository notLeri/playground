// Написать код, который будет выводить коды ответов на НТТР-запросы по URL адресам.
// Запросы должны осуществляться параллельно.
//  https://google.com - 200
//	https://dzen.ru/ - 200

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	urll = "https://google.com"
	url2 = "https://dzen.ru/"
)

func main() {
	ctx := context.Background()
	urls := []string{urll, url2}
	resultCh := make(chan string, len(urls))

	defaultClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}

	for _, url := range urls {
		go func(url string) {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			resp, err := defaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			resultCh <- fmt.Sprintf("%s - %d", url, resp.StatusCode)
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-resultCh)
	}
}