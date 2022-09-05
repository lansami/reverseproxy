package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lansami/reverse-proxy/httpproxy"
)

func main() {
	fmt.Println("Getting config file");
	configPath := flag.String("config", "./config.yaml", "Path to reverse proxy config");
	flag.Parse();
	
	config, err := ioutil.ReadFile(*configPath);
	if err != nil {
		log.Fatal("Could not read config file");
	}

	proxy := httpproxy.ReadProxyConfigFile(config);
	proxy.StartServer();
}