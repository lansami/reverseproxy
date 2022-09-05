package httpproxy

import (
	"testing"
)

func TestSetIsAlive(t *testing.T) {
	host := Host{};
	host.setIsAlive(true);

	if host.IsAlive() != true {
		t.Fail();
	}
}

func TestInitializeHost(t *testing.T) {
	host := Host{
		Address: "localhost",
		Port: 8888,
	}

	host.InitializeHost();

	if host.hostServer == nil {
		t.Fail();
	}
}