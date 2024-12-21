package consul

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/TylerAldrich814/MetaReviews/pkg/discovery"
	consul "github.com/hashicorp/consul/api"
)

// Registry defines a Consul-based service registry.
//
// Registry implements the MetaReviews/pkg/discovery/Registry interface
type Registry struct {
  client *consul.Client
}

// NewRegistry creates a new Consul-based service registry instance..
func NewRegistry(addr string)( *Registry, error ){
  config := consul.DefaultConfig()
  config.Address = addr

  client, err := consul.NewClient(config)
  if err != nil {
    return nil, err
  }
  return &Registry{ client }, nil
}

// Register creates a service revord in the consul-based registry.
func(r *Registry) Register(
  ctx         context.Context, 
  instanceID  string, 
  serviceName string, 
  hostPort string,
) error {
  parts := strings.Split(hostPort, ":")
  if len(parts) != 2 {
    return discovery.ErrInvalidHostPort
  }
  port, err := strconv.Atoi(parts[1])
  if err != nil {
    return err
  }

  return r.client.Agent().ServiceRegister(
    &consul.AgentServiceRegistration{
      Name    : serviceName,
      ID      : instanceID,
      Address : parts[0],
      Port    : port,
      Check   : &consul.AgentServiceCheck{
        CheckID : instanceID,
        TTL     : fmt.Sprintf(
          "%ds", 
          discovery.HealthCheckInterval,
        ),
      },
    },
  )
}

// Deregister removes a service revord in the consul-based registry.
func(r *Registry) Deregister(
  ctx         context.Context, 
  instanceID  string, 
  serviceName string,
) error {
  log.Printf("Deregistering %s instance.", serviceName)
  return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of active instances of the given service id.
func(r *Registry) ServiceAddresses(
  ctx         context.Context, 
  serviceName string,
)( []string, error) {
  entries, _, err := r.client.Health().Service(
    string(serviceName),
    "",
    true,
    nil,
  )
  if err != nil {
    return nil, err
  }
  if len(entries) == 0 {
    return nil, discovery.ErrNotFound
  }

  var res []string
  for _, e := range entries {
    res = append(
      res,
      fmt.Sprintf(
        "%s:%d",
        e.Service.Address,
        e.Service.Port,
      ),
    )
  }

  return res,nil
}

// ReportHealthyState is a mechanism for reporting the health status of the given instance to the registry.
func(r *Registry) ReportHealthyState(
  instanceID  string, 
  serviceName string,
) error {
  r.client.Agent().PassTTL(
    instanceID,
    "",
  )
  // TODO: Look more into UpdateTTL -- Find better examples of this newer API.
  return r.client.Agent().UpdateTTL(
    instanceID,
    "",
    "pass",
  )
}
