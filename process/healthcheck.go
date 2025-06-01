package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
)

// startHealthCheck starts the health check routine
func (sm *ServiceManager) startHealthCheck() {
	logger := sm.logger.Named("startHealthCheck")

	healthConfig := sm.config.HealthCheck.Unwrap()
	interval := healthConfig.Interval.UnwrapOr(30 * time.Second)

	sm.healthTicker = time.NewTicker(interval)

	logger.Info("starting health check", zap.Duration("interval", interval))

	go func() {
		for {
			select {
			case <-sm.ctx.Done():
				return
			case <-sm.healthTicker.C:
				sm.performHealthCheck()
			}
		}
	}()
}

// performHealthCheck performs a health check
func (sm *ServiceManager) performHealthCheck() {
	logger := sm.logger.Named("performHealthCheck")

	if !sm.config.HealthCheck.IsSome() {
		return
	}

	healthConfig := sm.config.HealthCheck.Unwrap()
	timeout := healthConfig.Timeout.UnwrapOr(10 * time.Second)
	retries := healthConfig.Retries.UnwrapOr(3)

	ctx, cancel := context.WithTimeout(sm.ctx, timeout)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		if attempt > 0 {
			logger.Debug("retrying health check", zap.Int("attempt", attempt+1))
			time.Sleep(time.Second)
		}

		switch healthConfig.Type {
		case "http":
			lastErr = sm.performHTTPHealthCheck(ctx, healthConfig)
		case "tcp":
			lastErr = sm.performTCPHealthCheck(ctx, healthConfig)
		case "cmd":
			lastErr = sm.performCommandHealthCheck(ctx, healthConfig)
		default:
			lastErr = fmt.Errorf("unknown health check type: %s", healthConfig.Type)
		}

		if lastErr == nil {
			logger.Debug("health check passed")
			return
		}

		logger.Warn("health check failed",
			zap.Error(lastErr),
			zap.Int("attempt", attempt+1))
	}

	logger.Error("health check failed after all retries", zap.Error(lastErr))
	sm.setState(StateHealthCheckFailed)
	sm.sendStatus(SignalError, StateHealthCheckFailed, lastErr)
}

// performHTTPHealthCheck performs HTTP health check
func (sm *ServiceManager) performHTTPHealthCheck(ctx context.Context, config *config.HealthCheckConfig) error {
	// Implementation would use http.Client with context
	// This is a placeholder
	return errors.New("HTTP health check not implemented")
}

// performTCPHealthCheck performs TCP health check
func (sm *ServiceManager) performTCPHealthCheck(ctx context.Context, config *config.HealthCheckConfig) error {
	// Implementation would use net.Dialer with context
	// This is a placeholder
	return fmt.Errorf("TCP health check not implemented")
}

// performCommandHealthCheck performs command-based health check
func (sm *ServiceManager) performCommandHealthCheck(ctx context.Context, config *config.HealthCheckConfig) error {
	logger := sm.logger.Named("performCommandHealthCheck")

	if len(config.Command) == 0 {
		return fmt.Errorf("no command specified for health check")
	}

	// Build command
	cmdParts := make([]string, 0, len(config.Command))
	for _, part := range config.Command {
		if part.IsSome() {
			cmdParts = append(cmdParts, *part.Unwrap())
		}
	}

	if len(cmdParts) == 0 {
		return fmt.Errorf("no valid command parts")
	}

	cmd := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)

	// Capture output if matcher is specified
	var stdoutBuf, stderrBuf *bytes.Buffer
	if config.ResultMatcher.IsSome() {
		stdoutBuf = &bytes.Buffer{}
		stderrBuf = &bytes.Buffer{}
		cmd.Stdout = stdoutBuf
		cmd.Stderr = stderrBuf
	}

	err := cmd.Run()
	// Check exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()

			// Check if exit code is acceptable
			for _, okCode := range config.OkExitCodes {
				if okCode.IsSome() && *okCode.Unwrap() == exitCode {
					err = nil
					break
				}
			}
		}
	}

	// Check output matcher if specified
	if err == nil && config.ResultMatcher.IsSome() {
		matcher := config.ResultMatcher.Unwrap()

		// Compile regex pattern
		regex, regexErr := regexp.Compile(*matcher)
		if regexErr != nil {
			logger.Error("invalid regex pattern in health check matcher",
				zap.Stringp("pattern", matcher),
				zap.Error(regexErr))
			return fmt.Errorf("invalid regex pattern: %w", regexErr)
		}

		// Combine stdout and stderr
		var combinedOutput strings.Builder
		if stdoutBuf != nil {
			combinedOutput.WriteString(stdoutBuf.String())
		}
		if stderrBuf != nil {
			if combinedOutput.Len() > 0 {
				combinedOutput.WriteString("\n")
			}
			combinedOutput.WriteString(stderrBuf.String())
		}

		output := combinedOutput.String()

		// Test regex against combined output
		if !regex.MatchString(output) {
			logger.Debug("health check output does not match pattern",
				zap.Stringp("pattern", matcher),
				zap.String("output", output))
			err = fmt.Errorf("output does not match expected pattern: %s", *matcher)
		} else {
			logger.Debug("health check output matches pattern",
				zap.Stringp("pattern", matcher))
		}
	}

	if err != nil {
		logger.Debug("command health check failed",
			zap.Error(err),
			zap.Strings("command", cmdParts))
	} else {
		logger.Debug("command health check passed",
			zap.Strings("command", cmdParts))
	}

	return err
}
