package config

import (
	"os"
)

// Service represents a single service definition in the system,
// including metadata, execution details, lifecycle, and dependencies.
type Template struct {
	// SourceDirectory where template files are stored.
	// This field is required.
	SourceDirectory string `mapstructure:"source,omitempty" yaml:"source" validate:"required"`

	// TargetDirectory where template files after compiling should be placed in.
	// defaults to [SourceDirectory].
	TargetDirectory string `mapstructure:"destination,omitempty" yaml:"destination" validate:"required"`

	// Extension of template files (will be removed after compile)
	// defaults to [DefaultTemplateExtension]
	Extension string `mapstructure:"extension,omitempty" yaml:"extension"`

	// Enabled specifies whether the template should be applied or not.
	// If false, the template directory will be ignored.
	// defaults to 'true'
	Enabled bool `mapstructure:"enabled,omitempty" yaml:"enabled" default:"true"`

	// OverWrite target file if exists.
	// defaults to [DefaultTemplateOverWrite]
	OverWrite bool `mapstructure:"overwrite,omitempty" yaml:"overwrite" default:"true"`

	// FileMod of target file that is created.
	// defaults to [DefaultTemplateFileMod]
	FileMod uint32 `mapstructure:"chmod,omitempty" yaml:"chmod" default:"0o644"`

	// DirMod of target directories that may get created during process.
	// defaults to [DefaultTemplateDirMod]
	DirMod uint32 `mapstructure:"dir_chmod,omitempty" yaml:"dir_chmod" default:"0o755"`

	// Fatal if true will deny the execution of the services.
	// defaults to [DefaultTemplateFatality]
	IsFatal bool `mapstructure:"is_fatal,omitempty" yaml:"is_fatal" default:"true"`
}

func (t *Template) GetFileMode() os.FileMode {
	return os.FileMode(t.FileMod)
}

func (t *Template) GetDirMode() os.FileMode {
	return os.FileMode(t.DirMod)
}
