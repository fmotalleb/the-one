package config

import (
	"time"
)

const (
	DefaultRestartDelayBegin = time.Second
	DefaultRestartDelayMax   = time.Second * 16
	DefaultServiceType       = OngoingService
	DefaultProcessCount      = 1
)

var (
	DefaultRestartOkCodes = []int{0}

	DefaultRestartConfig = RetryConfig{}
)
