// Package process encapsulates logic behind process management and policies for individual processes
package process

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/logging"
)

var log = logging.LazyLogger("process")

// NewServiceManager creates a new service manager.
func NewServiceManager(
	ctx context.Context,
	config *config.Service,
	controlCh chan ServiceMessage,
	statusCh chan ServiceMessage,
	std StdConfig,
) *ServiceManager {
	ctx, cancel := context.WithCancel(ctx)

	logger := log().
		Named("service").
		Named(config.GetName())

	return &ServiceManager{
		config:    config,
		instances: make([]*Instance, 0, config.GetProcessCount()),
		state:     StateStopped,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		controlCh: controlCh,
		statusCh:  statusCh,
		std:       std,
	}
}

// Start starts the service manager.
func (sm *ServiceManager) Start() error {
	logger := sm.logger.Named("Start")
	logger.Info("starting service manager", zap.String("service", sm.config.GetName()))

	sm.setState(StateStarting)

	// Start control signal listener
	go sm.handleControlSignals()

	// Start health check if configured
	if sm.config.HealthCheck.IsSome() {
		go sm.startHealthCheck()
	}

	// Start process instances
	if err := sm.startProcesses(); err != nil {
		sm.setState(StateError)
		return fmt.Errorf("failed to start processes: %w", err)
	}

	sm.setState(StateRunning)
	logger.Info("service manager started successfully")
	return nil
}

// Stop stops the service manager.
func (sm *ServiceManager) Stop() error {
	logger := sm.logger.Named("Stop")
	logger.Info("stopping service manager", zap.String("service", sm.config.GetName()))

	sm.setState(StateStopping)

	// Stop health check
	if sm.healthTicker != nil {
		sm.healthTicker.Stop()
	}

	// Stop all process instances
	if err := sm.stopProcesses(); err != nil {
		logger.Error("error stopping processes", zap.Error(err))
	}

	// Cancel context
	sm.cancel()

	sm.setState(StateStopped)
	logger.Info("service manager stopped")
	return nil
}

// handleControlSignals processes incoming control signals.
func (sm *ServiceManager) handleControlSignals() {
	logger := sm.logger.Named("handleControlSignals")

	for {
		select {
		case <-sm.ctx.Done():
			return
		case msg := <-sm.controlCh:
			logger.Debug("received control signal",
				zap.Int("signal", int(msg.Signal)),
				zap.String("service", msg.ServiceID))

			switch msg.Signal {
			case SignalStart:
				if err := sm.startProcesses(); err != nil {
					sm.sendStatus(SignalError, StateError, err)
				} else {
					sm.setState(StateRunning)
				}
			case SignalStop:
				if err := sm.stopProcesses(); err != nil {
					sm.sendStatus(SignalError, StateError, err)
				} else {
					sm.setState(StateStopped)
				}
			case SignalRestart:
				sm.restart()
			case SignalHealthCheck:
				go sm.performHealthCheck()
			}
		}
	}
}

// startProcesses starts all process instances.
func (sm *ServiceManager) startProcesses() error {
	logger := sm.logger.Named("startProcesses")
	processCount := sm.config.GetProcessCount()

	logger.Info("starting processes",
		zap.Int("count", processCount),
		zap.Stringp("executable", sm.config.Executable.Unwrap()))

	for i := 0; i < processCount; i++ {
		instance, err := sm.createProcessInstance(i)
		if err != nil {
			return fmt.Errorf("failed to create process instance %d: %w", i, err)
		}

		if err := sm.startProcessInstance(instance); err != nil {
			return fmt.Errorf("failed to start process instance %d: %w", i, err)
		}

		sm.instances = append(sm.instances, instance)
	}

	return nil
}

// createProcessInstance creates a new process instance.
func (sm *ServiceManager) createProcessInstance(id int) (*Instance, error) {
	logger := sm.logger.Named("createProcessInstance")

	// Prepare command
	executable := sm.config.Executable.Unwrap()
	args := make([]string, 0, len(sm.config.Arguments))

	for _, arg := range sm.config.Arguments {
		if arg.IsSome() {
			args = append(args, *arg.Unwrap())
		}
	}

	// #nosec G204 -- the command will be from config file so yeah its a variable
	cmd := exec.CommandContext(sm.ctx, *executable, args...)

	// Set working directory
	if sm.config.WorkingDir.IsSome() {
		cmd.Dir = *sm.config.WorkingDir.Unwrap()
	}

	cmd.Stdin = sm.std.In
	cmd.Stdout = sm.std.Out
	cmd.Stderr = sm.std.Err
	// Set environment variables
	env, err := sm.buildEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to build environment: %w", err)
	}
	cmd.Env = env

	// Set up process group for proper signal handling
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	instance := &Instance{
		ID:    id,
		Cmd:   cmd,
		State: StateStopped,
	}

	logger.Debug("created process instance",
		zap.Int("id", id),
		zap.Stringp("executable", executable),
		zap.Strings("args", args))

	return instance, nil
}

// startProcessInstance starts a single process instance
func (sm *ServiceManager) startProcessInstance(instance *Instance) error {
	logger := sm.logger.Named("startProcessInstance")

	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	logger.Info("starting process instance", zap.Int("id", instance.ID))

	instance.State = StateStarting
	instance.StartAt = time.Now()

	if err := instance.Cmd.Start(); err != nil {
		instance.State = StateError
		return fmt.Errorf("failed to start process: %w", err)
	}

	instance.State = StateRunning

	// Monitor process in background
	go sm.monitorProcessInstance(instance)

	logger.Info("process instance started",
		zap.Int("id", instance.ID),
		zap.Int("pid", instance.Cmd.Process.Pid))

	return nil
}

// monitorProcessInstance monitors a process instance
func (sm *ServiceManager) monitorProcessInstance(instance *Instance) {
	logger := sm.logger.Named("monitorProcessInstance")

	err := instance.Cmd.Wait()

	instance.mutex.Lock()
	instance.State = StateStopped
	instance.mutex.Unlock()

	if err != nil {
		logger.Error("process instance exited with error",
			zap.Int("id", instance.ID),
			zap.Error(err))
	} else {
		logger.Info("process instance exited normally", zap.Int("id", instance.ID))
	}
	// Check if restart is needed
	if sm.config.GetType() == config.OngoingService {
		go sm.restart()
	} else {
		sm.sendStatus(SignalError, StateError, err)
	}
}

// stopProcesses stops all process instances
func (sm *ServiceManager) stopProcesses() error {
	logger := sm.logger.Named("stopProcesses")

	var wg sync.WaitGroup
	errors := make(chan error, len(sm.instances))

	for _, instance := range sm.instances {
		wg.Add(1)
		go func(inst *Instance) {
			defer wg.Done()
			if err := sm.stopProcessInstance(inst); err != nil {
				errors <- err
			}
		}(instance)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	// Clear instances
	sm.instances = sm.instances[:0]

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping processes: %v", errs)
	}

	logger.Info("all processes stopped")
	return nil
}

// stopProcessInstance stops a single process instance
func (sm *ServiceManager) stopProcessInstance(instance *Instance) error {
	logger := sm.logger.Named("stopProcessInstance")

	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if instance.State == StateStopped {
		return nil
	}

	logger.Info("stopping process instance", zap.Int("id", instance.ID))

	instance.State = StateStopping

	// Send SIGTERM first
	if err := instance.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		logger.Warn("failed to send SIGTERM", zap.Error(err))
	}

	// Wait for graceful shutdown with timeout.
	timeout := defaultProcessTimeout
	if sm.config.Timeout.IsSome() {
		timeout = *sm.config.Timeout.Unwrap()
	}

	done := make(chan error, 1)
	go func() {
		done <- instance.Cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		logger.Warn("process did not exit gracefully, sending SIGKILL")
		if err := instance.Cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		<-done // Wait for process to actually exit
	case err := <-done:
		if err != nil {
			logger.Debug("process exited with error", zap.Error(err))
		}
	}

	instance.State = StateStopped
	logger.Info("process instance stopped", zap.Int("id", instance.ID))

	return nil
}

// buildEnvironment builds the environment variables for the process
func (sm *ServiceManager) buildEnvironment() ([]string, error) {
	logger := sm.logger.Named("buildEnvironment")

	env := make(map[string]string)

	// Start with current environment if passthrough is enabled
	if sm.config.EnvPassThru.UnwrapOr(false) {
		for _, e := range os.Environ() {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	}

	// Load environment files
	for _, envFile := range sm.config.EnvironmentFile {
		if envFile.IsSome() {
			if err := sm.loadEnvFile(*envFile.Unwrap(), env); err != nil {
				return nil, fmt.Errorf("failed to load env file %s: %w", *envFile.Unwrap(), err)
			}
		}
	}

	// Apply explicit environment variables
	for key, value := range sm.config.Environments {
		if value.IsSome() {
			env[key] = *value.Unwrap()
		} else {
			// Explicitly unset
			delete(env, key)
		}
	}

	// Convert to slice
	result := make([]string, 0, len(env))
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}

	logger.Debug("built environment", zap.Int("count", len(result)))
	return result, nil
}

// loadEnvFile loads environment variables from a file
func (sm *ServiceManager) loadEnvFile(filename string, env map[string]string) error {
	logger := sm.logger.Named("loadEnvFile")

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logger.Warn("invalid line in env file",
				zap.String("file", filename),
				zap.Int("line", lineNum),
				zap.String("content", line))
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	logger.Debug("loaded env file", zap.String("file", filename))
	return nil
}

// restart restarts the service.
func (sm *ServiceManager) restart() {
	logger := sm.logger.Named("restart")

	logger.Info("restarting service",
		zap.String("service", sm.config.GetName()),
		zap.Uint("attempt", sm.restartCount+1))

	sm.restartCount++
	sm.lastRestart = time.Now()

	// Calculate delay using the restart config
	restartConfig := sm.config.GetRestart()
	delay, shouldContinue := restartConfig.GetDelay(sm.restartCount - 1) // Use previous iteration for current delay

	if !shouldContinue {
		logger.Error("restart policy indicates should not continue after incrementing count")
		sm.setState(StateError)
		sm.sendStatus(SignalError, StateError, fmt.Errorf("restart limit exceeded"))
		return
	}

	if delay > 0 {
		logger.Info("waiting before restart",
			zap.Duration("delay", delay),
			zap.Uint("iteration", sm.restartCount))

		// Use context-aware sleep to allow cancellation during delay
		select {
		case <-sm.ctx.Done():
			logger.Warn("restart canceled during delay")
			return
		case <-time.After(delay):
			// Continue with restart
		}
	}

	logger.Info("executing restart", zap.Uint("attempt", sm.restartCount))

	// Stop current processes
	if err := sm.stopProcesses(); err != nil {
		logger.Error("failed to stop processes during restart", zap.Error(err))
		// Continue with restart attempt even if stop failed (possibly will happen)
	}

	// Start new processes
	if err := sm.startProcesses(); err != nil {
		logger.Error("failed to start processes during restart",
			zap.Error(err),
			zap.Uint("attempt", sm.restartCount))

		sm.setState(StateError)
		sm.sendStatus(SignalError, StateError, err)

		// Schedule another restart attempt if policy allows
		sm.restart()
	} else {
		logger.Info("service restarted successfully",
			zap.String("service", sm.config.GetName()),
			zap.Uint("attempt", sm.restartCount))

		sm.setState(StateRunning)
		sm.sendStatus(SignalStateChange, StateRunning, nil)
	}
}

// setState sets the service state and sends notification
func (sm *ServiceManager) setState(state ServiceState) {
	sm.mutex.Lock()
	oldState := sm.state
	sm.state = state
	sm.mutex.Unlock()

	if oldState != state {
		sm.sendStatus(SignalStateChange, state, nil)
	}
}

// sendStatus sends a status message
func (sm *ServiceManager) sendStatus(signal ServiceSignal, state ServiceState, err error) {
	select {
	case sm.statusCh <- ServiceMessage{
		Signal:    signal,
		ServiceID: sm.config.GetName(),
		State:     state,
		Error:     err,
	}:
	default:
		// Don't block if channel is full
	}
}

// GetState returns the current service state
func (sm *ServiceManager) GetState() ServiceState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.state
}

// GetInstances returns the current process instances
func (sm *ServiceManager) GetInstances() []*Instance {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	instances := make([]*Instance, len(sm.instances))
	copy(instances, sm.instances)
	return instances
}
