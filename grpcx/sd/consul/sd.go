package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/zhangliang-zl/reskit/grpcx/sd"
	"google.golang.org/grpc/resolver"
)

type discovery struct {
	ca string
}

func (d *discovery) Register(svc sd.Service) error {
	agent, err := getAgent(d.ca)
	if err != nil {
		return err
	}

	deregister := svc.HealthCheckInterval * 5

	reg := &api.AgentServiceRegistration{
		ID:      buildID(svc),
		Name:    svc.Name,
		Tags:    svc.Tag,
		Port:    svc.Port,
		Address: svc.IP,
		Check: &api.AgentServiceCheck{
			Interval:                       svc.HealthCheckInterval.String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", svc.IP, svc.Port, svc.Name), //grpc 支持，执行健康检查的地址，svc 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: deregister.String(),                                 // 注销时间，相当于过期时间
		},
	}

	return agent.ServiceRegister(reg)
}

func (d *discovery) DeRegister(svc sd.Service) error {
	agent, err := getAgent(d.ca)
	if err != nil {
		return err
	}
	id := buildID(svc)
	return agent.ServiceDeregister(id)
}

func (*discovery) ResolverBuilder() resolver.Builder {
	return &resolverBuilder{}
}

func getAgent(ca string) (*api.Agent, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = ca
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return client.Agent(), nil
}

func buildID(svc sd.Service) string {
	return fmt.Sprintf("%v-%v-%v", svc.Name, svc.IP, svc.Port)
}

func New(address string) sd.ServiceDiscovery {
	return &discovery{
		ca: address,
	}
}
