package config

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
	"net"
)

type ProxyType string

const (
	Nginx   ProxyType = "Nginx"
	HAProxy ProxyType = "HAProxy"
)

type ProxyAddress net.IPAddr

type ProxyAuthentication string

const (
	None     ProxyAuthentication = "none"
	Password ProxyAuthentication = "Password"
)

type NetDevice string

type addressPool struct {
	Name          string
	CIDR          []*net.IPNet
	AvoidBuggyIPs bool  `yaml:"avoid-buggy-ips"`
	AutoAssign    *bool `yaml:"auto-assign"`
	Proxy         string
}

//internal Type
type Pool struct {
	CIDR          []*net.IPNet
	AvoidBuggyIPs bool `yaml:"avoid-buggy-ips"`
	AutoAssign    bool `yaml:"auto-assign"`
	Proxy         *Proxy
}

type Proxy struct {
	Name                string
	ProxyType           ProxyType
	ProxyAddress        ProxyAddress
	ProxyAuthentication ProxyAuthentication
	NetDevice           NetDevice
}

type Config struct {
	Pool map[string]*Pool
}

type configFile struct {
	Pool    []addressPool `yaml:"address-pools"`
	Proxies []Proxy       `yaml:"proxies"`
}

func parseProxy(p Proxy) (*Proxy, error) {
	if p.Name == "" {
		return nil, errors.New("missing proxy name")
	}
	res := &Proxy{}
	res.Name = p.Name
	res.ProxyAddress = p.ProxyAddress

	switch p.ProxyAuthentication {
	case None:
		res.ProxyAuthentication = p.ProxyAuthentication
	case Password:
		res.ProxyAuthentication = p.ProxyAuthentication
	default:
		return nil, fmt.Errorf("unknown Proxy Authentication %s", p.ProxyAuthentication)
	}

	if p.ProxyType == "" {
		return nil, errors.New("missing Proxy type")
	}
	res.ProxyType = p.ProxyType
	switch p.ProxyType {
	case HAProxy:
		res.ProxyType = p.ProxyType
	case Nginx:
		res.ProxyType = p.ProxyType
	default:
		return nil, fmt.Errorf("unknown Proxy Type %s", p.ProxyType)

	}
	//TODO add Pools
	//TODO Validate the Rest
	return res, nil
}

func parsePool(p addressPool) (*Pool, error) {
	res := &Pool{}
	//TODO Parse Pools
	return res, nil
}

func Parse(bs []byte, validate Validate) (*Config, error) {
	var raw configFile
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return nil, fmt.Errorf("could not parse config: %s", err)
	}
	err := validate(&raw)
	if err != nil {
		return nil, err
	}

	proxies := map[string]*Proxy{}

	for i, proxy := range raw.Proxies {
		parsedProxy, err := parseProxy(proxy)
		if err != nil {
			fmt.Errorf("parsing proxy profile #%d: %s", i+1, err)
		}
		if _, ok := proxies[proxy.Name]; ok {
			return nil, fmt.Errorf("found duplicate %s", err)
		}
		proxies[proxy.Name] = parsedProxy
	}

	cfg := &Config{
		Pool: map[string]*Pool{},
	}

	for i, ap := range raw.Pool {
		pool, err := parsePool(ap)
		if err != nil {
			fmt.Errorf("parsing pool profile #%d: %s", i+1, err)
		}
		if _, ok := cfg.Pool[ap.Name]; ok {
			return nil, fmt.Errorf("found duplicate %s", err)
		}
		cfg.Pool[ap.Name] = pool
		cfg.Pool[ap.Name].Proxy = proxies[pool.Proxy.Name]

	}
	return cfg, nil
}
