package httpproxy

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy ReverseProxy `yaml:"proxy"`
} 

type ReverseProxy struct {
	Listen Listen `yaml:"listen"`
	Services []Service `yaml:"services"`
	Server *http.Server
}

type Listen struct {
	Address string `yaml:"address"`
	Port int `yaml:"port"`
}

type Service struct {
	Name string `yaml:"name"`
	Domain string `yaml:"domain"`
	LbPolicy string `yaml:"lbPolicy"`
	Timeout string `yaml:"timeout"`
	Retries int `yaml:"retries"`
	Hosts []Host `yaml:"hosts"`
	HealthCheckTimeout string `yaml:"healthCheckTimeout"`
	requestTimeout time.Duration
	nextHostIndexToTry int
}

type Host struct {
	Address string `yaml:"address"`
	Port int `yaml:"port"`
	Status Status
	hostServer *httputil.ReverseProxy
}

type Status struct {
	IsAlive bool
	mu sync.Mutex
}

func (rp *ReverseProxy) ServerHandler(wr http.ResponseWriter, req *http.Request) {
	host := req.Host;
	//Plan:
	// 1.Find upstream service
	for i:=0; i<len(rp.Services); i++ {
		service := rp.Services[i];

		if service.Domain == host {
			log.Printf("Found domain, redirecting request")
			service.HandleRequest(wr, req);
			break;
		}
	}
	// 2.Use that upstream service to process the request
	// 3.Send the response back to the client
}

func writeResponseWithMessage(wr http.ResponseWriter, status int, msg string) {
	wr.WriteHeader(status);
}

func (p *ReverseProxy) StartServer() {
	proxyAddress := p.Listen.Address;
	proxyPort := p.Listen.Port;

	var addr = flag.String("addr", proxyAddress + ":" + strconv.Itoa(proxyPort), "The address of the application.")
	flag.Parse()

	serverHandler := http.HandlerFunc(p.ServerHandler);

	p.Server = &http.Server{
		Addr: *addr,
		Handler: serverHandler,
	} 
	
	for i:= 0; i<len(p.Services); i++ {
		service := p.Services[i];
		service.InitializeUpstreamService();
	}

	log.Println("Starting proxy server on", *addr)
	if err := p.Server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	} else {
		log.Println("Server started")
	}
	
}

func ReadProxyConfigFile(config []byte) *ReverseProxy{
	proxyConfig := &Config{};
	
	var err = yaml.Unmarshal(config, proxyConfig);

	if err != nil {
		panic(err);
	}

	return &proxyConfig.Proxy;
}