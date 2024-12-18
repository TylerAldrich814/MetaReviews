package memory

import (
	"context"
	"sync"
	"time"

  "github.com/TylerAldrich814/MetaMovies/pkg/discovery"
)

type InMemoryRegistry = map[string]map[string]*serviceInstance

type serviceInstance struct {
  hostPort   string
  lastActive time.Time
}


// Registry implements the MetaMovies/pkg/discovery/Registry interface
// Registry holds mapped data for all Registered Service Instances.
// Used for testing or when this applicaiton is running on a simple 
// Server, where a heavy-duty discovery layer isn't quite needed.
type Registry struct {
  sync.RWMutex
  serviceAddrs InMemoryRegistry
}

// NewRegistry creates a new in-memory registry.
func NewRegistry() *Registry {
  return &Registry{
    serviceAddrs: InMemoryRegistry{},//map[ServiceName]map[InstanceID]*serviceInstance{},
  }
}

// Register creates a service record in the registry
func(r *Registry) Register(
  ctx         context.Context,
  instanceID  string,
  serviceName string,
  hostPort    string,
) error {
  r.Lock()
  defer r.Unlock()

  if _, ok := r.serviceAddrs[serviceName]; !ok {
    r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
  }
  r.serviceAddrs[serviceName][instanceID] = &serviceInstance{
    hostPort   : hostPort,
    lastActive : time.Now(),
  }

  return nil
}

// Deregister removes a service record from the registry.
func(r *Registry) Deregister(
  ctx         context.Context,
  instanceID  string,
  serviceName string,
) error {
  r.Lock()
  defer r.Unlock()

  if _, ok := r.serviceAddrs[serviceName]; !ok {
    return nil
  }

  delete(r.serviceAddrs[serviceName], instanceID)

  return nil
}

// ServiceAddresses iterates over the provided serviceName, retuning an 
// array of service addresses if any are found within the registry.
// Only returns instances that reported as healthy within the last 
// 'HealthCheckInterval' seconds
func(r *Registry) ServiceAddresses(
  ctx         context.Context, 
  serviceName string,
)( []string, error) {
  r.Lock()
  defer r.Unlock()

  if len(r.serviceAddrs[serviceName]) == 0 {
    return nil, discovery.ErrNotFound
  }
  var res []string
  for _, i := range r.serviceAddrs[serviceName] {
    if i.lastActive.Before(
      time.
        Now().
        Add(-discovery.HealthCheckInterval),
    ){
      continue
    }
    res = append(res, i.hostPort)
  }

  return res, nil
}

// ReportHealthyState Reports the health state of a given Registry.
func(r *Registry) ReportHealthyState(
  instanceID  string, 
  serviceName string,
) error {
  r.Lock()
  defer r.Unlock()

  if _, ok := r.serviceAddrs[serviceName]; !ok {
    return discovery.ErrServiceNotRegistered
  }
  if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
    return discovery.ErrInstanceNotRegistered
  }

  r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
  return nil
}
