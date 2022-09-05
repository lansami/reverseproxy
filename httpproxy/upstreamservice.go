package httpproxy

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (s *Service) InitializeUpstreamService() {
	s.nextHostIndexToTry = 0;

	if s.Retries <= 0 {
		s.Retries = 3;
	}

	if len(s.Timeout) < 2 {
		s.requestTimeout, _ = time.ParseDuration("10s");
	} else {
		s.requestTimeout, _ = time.ParseDuration(s.Timeout);
	}

	for i := 0; i < len(s.Hosts); i++ {
		host := &s.Hosts[i];
		host.InitializeHost();

		log.Printf("Host %s:%s is alive(%t)", host.Address, strconv.Itoa(host.Port), host.IsAlive())

		if !host.IsAlive() {
			for j := 0; j < s.Retries; j++ {
				host.InitializeHost();
				if host.IsAlive() {
					break;
				}
			} 
		}
	}

	go s.startCheckHostsHealthPeriodically();
}

func (s *Service) HandleRequest(wr http.ResponseWriter, req *http.Request) {
	host := s.getHost();

	if host != nil {
		host.ProcessRequest(wr, req);
	} else {
		log.Println("Could not find available host");
		wr.WriteHeader(500);
	}
}

func (s *Service) getHost() *Host {
	switch(s.LbPolicy) {
	case "RANDOM":
		return s.getRandomHost();
	case "ROUND_ROBIN":
		return s.getRoundRobinHost();
	default:
		return s.getRandomHost();
	}
}

func (s *Service) getRandomHost() *Host{
	hostsCount := len(s.Hosts);
	if hostsCount > 0 {
		tries := 0;
		for ok := true ; ok; ok = tries < hostsCount*2 {
			index := rand.Intn(hostsCount);
			host := &s.Hosts[index];
			if !host.IsAlive() {
				tries++;
				continue;
			}

			return host;
		}
	}

	return nil;
}

func (s *Service) getRoundRobinHost() *Host{
	hostCount := len(s.Hosts);
	host := &s.Hosts[s.nextHostIndexToTry%hostCount];
	s.nextHostIndexToTry++;
	tries := 0;
	for ok := false; ok; ok = !host.IsAlive() || tries < s.Retries {
		host = &s.Hosts[s.nextHostIndexToTry%hostCount];
		s.nextHostIndexToTry++;
		tries++;
	}
	if host.IsAlive() {
		return host;
	} else {
		return nil;
	}
}

func (s *Service) startCheckHostsHealthPeriodically() {
	c := make(chan string);
	go s.checkHostsHealth(c);
	for message := range c {
		go func(msg string) {
			healthCheckTimeoutString := s.HealthCheckTimeout;
			if len(healthCheckTimeoutString) < 2 {
				healthCheckTimeoutString = "10s";
			}
			healthCheckTimeout, _ := time.ParseDuration(healthCheckTimeoutString);
			time.Sleep(healthCheckTimeout * time.Second);
			s.checkHostsHealth(c);
			log.Println(msg);
		}(message)
	}
}

func (s *Service) checkHostsHealth(c chan string) {
	for i:=0; i < len(s.Hosts); i++ {
		host := &s.Hosts[i];
		host.checkHealth();
		c <- "HOST_CHECKED"
	}
}
