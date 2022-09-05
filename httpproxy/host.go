package httpproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func (h *Host) setIsAlive(isAlive bool) {
	h.Status.mu.Lock();
	h.Status.IsAlive = isAlive;
	h.Status.mu.Unlock();
}

func (h *Host) IsAlive() bool {
	isAlive := false;
	h.Status.mu.Lock();
	isAlive = h.Status.IsAlive;
	h.Status.mu.Unlock();

	return isAlive;
}

func (h *Host) InitializeHost() {
	targetUrl := url.URL{
		Scheme: "http",
		Host: fmt.Sprintf("%s:%d", h.Address, h.Port),
	}
	h.hostServer = httputil.NewSingleHostReverseProxy(&targetUrl);
	h.hostServer.ErrorHandler = func(w http.ResponseWriter,r *http.Request, e error) {
		log.Printf("Error while sending request to %s", &targetUrl);
		h.setIsAlive(false);
	}
	h.setIsAlive(true);
}

func (h *Host) ProcessRequest(wr http.ResponseWriter, req *http.Request) {
	log.Printf("Processing request on host: %s:%s", h.Address, strconv.Itoa(h.Port));
	h.hostServer.ServeHTTP(wr, req);
}

func (h *Host) checkHealth() {
	hostAddress := fmt.Sprintf("http://%v:%d",h.Address, h.Port);
	_, err := http.Get(hostAddress)

	if err != nil {
		log.Printf("Host %v is not alive", hostAddress)
		log.Print(err);
		h.setIsAlive(false);
	} else {
		log.Printf("Host %v is alive", hostAddress)
		h.setIsAlive(true);
	}
}