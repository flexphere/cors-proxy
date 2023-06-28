package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var port string
var target string

func init() {
	flag.StringVar(&port, "p", "8080", "port to serve on")
	flag.StringVar(&target, "t", "", "proxy target url eg:http://example.com")
	flag.Parse()

	if target == "" {
		flag.Usage()
		log.Fatalln("ERR: target url must be set")
	}
}

func main() {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal("Failed to parse target URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	proxy.Director = func(req *http.Request) {
		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.Header.Set("Connection", "keep-alive")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("[%d] %s\n", resp.StatusCode, resp.Request.URL.String())
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Headers", "*")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Allow-Methods", "*")
		return nil
	}

	log.Printf("Redirecting http://localhost:%s -> %s\n", port, target)
	err = http.ListenAndServe(":"+port, proxy)
	if err != nil {
		log.Fatal("Failed to start proxy server:", err)
	}
}
