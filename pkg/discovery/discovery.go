package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// HealthCheckInterval - a constant that determines how long each 
// Service Instance Health Check interval should last.
const HealthCheckInterval = 5 * time.Second

// // string - The ID of specific service instance.
// type InstanceID string
// // string - The Name of specific service instance.
// type ServiceName string

// Registry defines a service registry.
type Registry interface {
  // Register creates a service instance record in the registry,
  Register(
    ctx         context.Context, 
    instanceID  string, 
    serviceName string, 
    hostPort    string,
  ) error

  // Deregister removes a service instance record from the registry.
  Deregister(
    ctx         context.Context, 
    instanceID  string, 
    serviceName string,
  ) error

  // ServiceAddresses returns the list of addresses of active instances of the given service id.
  ServiceAddresses(
    ctx         context.Context, 
    serviceName string,
  )( []string, error)

  // ReportHealthyState is a mechanism for reporting the health status of the given instance to the registry.
  ReportHealthyState(
    instanceID  string, 
    serviceName string,
  ) error
}

// ErrNotFound is returned when no service addresses are found.
var (
  ErrNotFound              = errors.New("no service address found")
  ErrServiceNotRegistered  = errors.New("the provided service is not registered")
  ErrInstanceNotRegistered = errors.New("the provided instance is not registered")
  ErrInvalidHostPort       = errors.New("the provided hostPort address is incorrectly formatted")
)

// Generatestring generates a pseudo-random service instance identifier, using a service name
// suffixed by a dash and a random number.
func GenerateInstanceID(serviceName string) string {
  return fmt.Sprintf(
    "%s-%d",
    serviceName,
    rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
  )
}
