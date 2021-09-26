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

type ServiceDiscovery interface {
	Register(svc Service) error
	DeRegister(svc Service) error
	ResolverBuilder() resolver.Builder
}