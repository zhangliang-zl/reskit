package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"regexp"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")
	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z0-9_]+)$")
)

type resolverBuilder struct {
}

func (b *resolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	host, port, name, err := b.parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		return nil, err
	}

	cr := &resolverImpl{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	go cr.watcher()
	return cr, nil
}

func (b *resolverBuilder) parseTarget(target string) (host, port, name string, err error) {

	if target == "" {
		return "", "", "", errMissingAddr
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	if port == "" {
		port = "8500"
	}

	return host, port, name, nil
}

func (b *resolverBuilder) Scheme() string {
	return "consul"
}

type resolverImpl struct {
	address              string
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}

func (r *resolverImpl) watcher() error {
	config := api.DefaultConfig()
	config.Address = r.address
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	for {
		services, meta, err := client.Health().Service(r.name, r.name, true, &api.QueryOptions{WaitIndex: r.lastIndex})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v\n", err)
			return err
		}

		r.lastIndex = meta.LastIndex
		var addresses []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			addresses = append(addresses, resolver.Address{Addr: addr})
		}
		state := resolver.State{
			Addresses:     addresses,
		}
		r.cc.UpdateState(state)
	}
}

func (r *resolverImpl) ResolveNow(opt resolver.ResolveNowOptions) {
}

func (r *resolverImpl) Close() {
}
