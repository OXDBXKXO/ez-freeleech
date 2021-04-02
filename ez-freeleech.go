package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"

	golog "github.com/mandala/go-log"
)

var logger = golog.New(os.Stdout).WithColor().WithDebug()
var locale = "en"

func getReadableBytes(bytes int64, locale string) string {
	en := []string{"B", "kB", "MB", "GB", "TB"}
	fr := []string{"o", "ko", "Mo", "Go", "To"}

	var pow int64 = 4
	for pow > 0 && bytes <= int64(math.Pow(1000, float64(pow))) {
		pow--
	}

	if locale == "en" {
		return fmt.Sprintf("%d %s", bytes/int64(math.Pow(1000, float64(pow))), en[pow])
	} else {
		return fmt.Sprintf("%d %s", bytes/int64(math.Pow(1000, float64(pow))), fr[pow])
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	event, queryContainsEvent := query["event"]
	downloaded, queryContainsDownloaded := query["downloaded"]

	if queryContainsEvent {
		if queryContainsDownloaded && len(downloaded) > 0 && downloaded[0] != "0" {
			i, err := strconv.ParseInt(downloaded[0], 10, 64)
			if err != nil {
				logger.Errorf("Could not parse %s to Int64\n", downloaded[0])
				i = -1
			}
			logger.Infof("%s %s tried to report %s\n", event, req.Host, getReadableBytes(i, locale))

		} else {
			logger.Infof("%s %s", event, req.Host)
		}
	}

	freeleechURI := req.RequestURI
	downloadedRegex, err := regexp.Compile(`downloaded=[0-9]+?&`)
	if err != nil {
		logger.Warn("Could not find 'downloaded' key in query string")
	} else {
		freeleechURI = downloadedRegex.ReplaceAllString(req.RequestURI, "downloaded=0&")
	}

	freeleech, err := http.NewRequest("GET", freeleechURI, nil)
	if err != nil {
		logger.Error("Could not build freeleech request, aborting request")
		logger.Info("This will prevent download report to the tracker but you will not receive peer list")
		return
	}

	freeleech.Header.Set("User-Agent", "qBittorrent/4.3.3")
	freeleech.Header.Set("Accept-Encoding", "gzip")
	freeleech.Header.Set("Connection", "Close")

	resp, err := http.DefaultTransport.RoundTrip(freeleech)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8888, "Port to run ez-freeleech on")
	flag.StringVar(&locale, "locale", "en", "Locale for data representation (en|fr)")
	flag.Parse()

	if locale != "fr" && locale != "en" {
		logger.Fatal("locale cannot be set to '" + locale + "'. Only 'en' and 'fr' are allowed values")
	}

	logger.Info(fmt.Sprintf("Started ez-freeleech on port %d", port))

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
		})}

	logger.Fatal(server.ListenAndServe())
}
