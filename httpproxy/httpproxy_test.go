package httpproxy

import (
	"flag"
	"io/ioutil"
	"testing"
)

func TestReadProxyConfigFile(t *testing.T) {
	configPath := flag.String("config", "./config_test.yaml", "Path to reverse proxy config");
	flag.Parse();
	config, err := ioutil.ReadFile(*configPath);
	if err != nil {
		t.Error(err);
	}
	proxy := ReadProxyConfigFile(config);

	if proxy.Listen.Address == "" || proxy.Listen.Port <= 0{
		t.Fail();
	}
}