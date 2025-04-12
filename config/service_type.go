package config

type ServiceType string

const (
	OngoingService ServiceType = "normal"
	OneShotService ServiceType = "one-shot"
	MultiShot      ServiceType = "multi-shot"
)

func (s ServiceType) IsValid() bool {
	switch s {
	case OngoingService, OneShotService, MultiShot:
		return true
	default:
		return false
	}
}
