package config

import (
	"time"
)

const (
	DefaultRestartDelayBegin = time.Second
	DefaultRestartDelayMax   = time.Second * 16
	DefaultServiceType       = OngoingService
	DefaultProcessCount      = 1

	DefaultTemplateExtension = ".template"
	DefaultTemplateOverWrite = true
	DefaultTemplateFileMod   = 0o644
	DefaultTemplateDirMod    = 0o755
	DefaultTemplateFatality  = true
)

var (
	DefaultRestartOkCodes = []int{0}

	DefaultRestartConfig = RetryConfig{}
)
