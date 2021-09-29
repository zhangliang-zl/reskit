package sd

import (
	"google.golang.org/grpc/resolver"
	"time"
)

type Service struct {
	IP                  string
	Port                int
	Tag                 []string
	Name                string
	HealthCheckInterval time.Duration
}

type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpcx://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"endpoints"`
}

type Register interface {
	Register() error
	DeRegister() error
}

type ServiceDiscovery interface {
	Register(svc Service) error
	DeRegister(svc Service) error
	ResolverBuilder() resolver.Builder
}
