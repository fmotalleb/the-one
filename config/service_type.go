package config

import "fmt"

type ServiceType string

const (
	// OngoingService is a kind of service that runs continuously and can be restarted.
	OngoingService ServiceType = "normal"

	// OneShotService is a kind of service that runs once will be executed again if exit code was non-zero value.
	OneShotService ServiceType = "one-shot"

	// MultiShotService is a kind of service that runs once per dependency and will be executed again if exit code was non-zero value.
	MultiShotService ServiceType = "multi-shot"
)

func (s ServiceType) IsValid() bool {
	switch s {
	case OngoingService, OneShotService, MultiShotService:
		return true
	default:
		return false
	}
}

func (s *ServiceType) Parse(a any) (ServiceType, error) {
	if str, ok := a.(string); ok {
		*s = ServiceType(str)
		if !s.IsValid() {
			return "", fmt.Errorf("invalid service type: %s", str)
		}
		return *s, nil
	}
	return "", fmt.Errorf("invalid service type: %v", a)
}
