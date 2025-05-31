package config

import (
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

const (
	DefaultRestartDelayBegin = time.Second
	DefaultRestartDelayMax   = time.Second * 16
	DefaultServiceType       = OngoingService
	DefaultProcessCount      = 1

	restartDelayPowerBase         = 2
	restartMaxCalculableIteration = 10
	restartAbsoluteMax            = uint(1000000)

	DefaultTemplateExtension = ".template"
	DefaultTemplateOverWrite = true
	DefaultTemplateFileMod   = 0o644
	DefaultTemplateDirMod    = 0o755
	DefaultTemplateFatality  = true
)

var (
	DefaultRestartOkCodes = []int{0}

	DefaultRestartConfig = RestartConfig{
		Count:    option.Optional[uint]{},
		Delay:    option.Optional[time.Duration]{},
		DelayMax: option.Optional[time.Duration]{},
	}
)
