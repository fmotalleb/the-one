package process

import (
	"context"
	"io"
	"os/exec"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
)

// ServiceState represents the current state of a service.
type ServiceState int

const (
	StateStopped ServiceState = iota
	StateStarting
	StateRunning
	StateStopping
	StateError
	StateHealthCheckFailed
)

// ServiceSignal represents signals sent between service and controller.
type ServiceSignal int

const (
	SignalStart ServiceSignal = iota
	SignalStop
	SignalRestart
	SignalHealthCheck
	SignalStateChange
	SignalError
)

// ServiceMessage represents a message sent via signal channels.
type ServiceMessage struct {
	Signal    ServiceSignal
	ServiceID string
	State     ServiceState
	Error     error
	Data      interface{}
}

// Instance represents a single process instance.
type Instance struct {
	ID      int
	Cmd     *exec.Cmd
	State   ServiceState
	StartAt time.Time
	mutex   sync.RWMutex
}

// ServiceManager manages a single service with multiple process instances.
type ServiceManager struct {
	config    *config.Service
	instances []*Instance
	state     ServiceState
	logger    *zap.Logger
	ctx       context.Context
	cancel    context.CancelFunc

	// Signal channels
	controlCh chan ServiceMessage // Receives control signals
	statusCh  chan ServiceMessage // Sends status updates

	// Health check
	healthTicker *time.Ticker

	// Restart tracking
	restartCount uint
	lastRestart  time.Time

	mutex sync.RWMutex

	std StdConfig
}

type StdConfig struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
