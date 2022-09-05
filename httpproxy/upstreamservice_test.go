package httpproxy

import (
	"testing"
)

func TestGetRoundRobinHost(t *testing.T) {
	s := Service{
		Hosts: []Host{},
		nextHostIndexToTry: 1,
	}
	s.Hosts = append(s.Hosts, Host{
		Address: "host1",
	});
	s.Hosts = append(s.Hosts, Host{
		Address: "host2",
	});
	s.Hosts = append(s.Hosts, Host{
		Address: "host3",
	});

	for i:=0;i<3;i++ {
		s.Hosts[i].setIsAlive(true);
	}

	host := s.getRoundRobinHost();

	if host.Address != "host2" {
		t.Fail();
	}
}

func TestGetRandomHostGetNullWhenAllHostDown(t *testing.T) {
	s := Service{
		Hosts: []Host{},
		nextHostIndexToTry: 1,
	}
	s.Hosts = append(s.Hosts, Host{
		Address: "host1",
	});
	s.Hosts = append(s.Hosts, Host{
		Address: "host2",
	});
	s.Hosts = append(s.Hosts, Host{
		Address: "host3",
	});

	host := s.getRandomHost();

	if host != nil {
		t.Fail();
	}
}

func TestGetRoundRobinHostNullWhenAllHostDown(t *testing.T) {
	s := Service{
		Hosts: []Host{},
		nextHostIndexToTry: 1,
	}
	s.Hosts = append(s.Hosts, Host{
		Address: "host1",
	});
	s.Hosts = append(s.Hosts, Host{
		Address: "host2",
	});

	host := s.getRoundRobinHost();

	if host != nil {
		t.Fail();
	}
}