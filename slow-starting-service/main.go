package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	startupDelayMax := flag.Int("delay-max", 60, "Maximum number of seconds to wait before accepting HTTP requests")
	startupDelayMin := flag.Int("delay-min", 10, "Minimum number of seconds to wait before accepting HTTP requests")
	listen := flag.String("listen", ":8080", "Interface to listen on")
	failurePercent := flag.Int("failure-percent", 1, "Probability of a health check failing that should fail")
	check := flag.Bool("check", false, "If true, check health of service and exit")

	flag.Parse()
	rand.Seed(time.Now().Unix())
	if *check == true {
		healthcheck(*listen)
	}

	http.HandleFunc("/health", health(*failurePercent, healthy, unhealthy))

	if *startupDelayMax < *startupDelayMin {
		tmp := *startupDelayMax
		*startupDelayMin = *startupDelayMax
		*startupDelayMax = tmp
	}

	waitForN := *startupDelayMax - *startupDelayMin

	duration := time.Duration(*startupDelayMin+rand.Intn(waitForN)) * time.Second
	log.Printf("Waiting for %s", duration)
	time.Sleep(duration)
	log.Printf("Starting HTTP server on %s", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen %s: %s\n", *listen, err)
	}
}

func health(failurePercent int, healthy, unhealthy http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		randomNumber := rand.Intn(100)
		if randomNumber <= failurePercent {
			unhealthy(w, req)
		} else {
			healthy(w, req)
		}
	}
}

func unhealthy(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"status": "error"}`)
	fmt.Fprintf(w, "\n")
}
func healthy(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, `{"status": "healthy"}`)
	fmt.Fprintf(w, "\n")
}

func healthcheck(listen string) {
	healthCheckURL := url.URL{
		Scheme: "http",
		Host:   listen,
		Path:   "/health",
	}
	response, err := http.Get(healthCheckURL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: http.Get: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	io.Copy(os.Stdout, response.Body)
	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "healthcheck: response status: %s\n", response.Status)
		os.Exit(1)
	}

	os.Exit(0)
}
