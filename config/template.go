package config

import (
	"os"

	"github.com/fmotalleb/the-one/types/option"
)

// Service represents a single service definition in the system,
// including metadata, execution details, lifecycle, and dependencies.
type Template struct {
	// SourceDirectory where template files are stored.
	// This field is required.
	SourceDirectory option.Some[string] `mapstructure:"source,omitempty" yaml:"source"`

	// TargetDirectory where template files after compiling should be placed in.
	// defaults to [SourceDirectory].
	TargetDirectory option.OptionalT[string] `mapstructure:"destination,omitempty" yaml:"destination"`

	// Extension of template files (will be removed after compile)
	// defaults to [DefaultTemplateExtension]
	Extension option.OptionalT[string] `mapstructure:"extension,omitempty" yaml:"extension"`

	// Enabled specifies whether the template should be applied or not.
	// If false, the template directory will be ignored.
	// defaults to 'true'
	Enabled option.OptionalT[bool] `mapstructure:"enabled,omitempty" yaml:"enabled"`

	// OverWrite target file if exists.
	// defaults to [DefaultTemplateOverWrite]
	OverWrite option.OptionalT[bool] `mapstructure:"overwrite,omitempty" yaml:"overwrite"`

	// FileMod of target file that is created.
	// defaults to [DefaultTemplateFileMod]
	FileMod option.OptionalT[uint32] `mapstructure:"chmod,omitempty" yaml:"chmod"`

	// DirMod of target directories that may get created during process.
	// defaults to [DefaultTemplateDirMod]
	DirMod option.OptionalT[uint32] `mapstructure:"dir_chmod,omitempty" yaml:"dir_chmod"`

	// Fatal if true will deny the execution of the services.
	// defaults to [DefaultTemplateFatality]
	IsFatal option.OptionalT[bool] `mapstructure:"is_fatal,omitempty" yaml:"is_fatal"`
}

func (t *Template) GetSourceDirectory() string {
	dir := *t.SourceDirectory.Unwrap()
	return dir
}

func (t *Template) GetFinalDirectory() string {
	return *t.TargetDirectory.UnwrapOr(t.GetSourceDirectory())
}

func (t *Template) GetExtension() string {
	return *t.Extension.UnwrapOr(DefaultTemplateExtension)
}

func (t *Template) IsEnabled() bool {
	return *t.Enabled.UnwrapOr(true)
}

func (t *Template) ShouldOverwrite() bool {
	return *t.OverWrite.UnwrapOr(DefaultTemplateOverWrite)
}

func (t *Template) GetFileMode() os.FileMode {
	return os.FileMode(*t.FileMod.UnwrapOr(DefaultTemplateFileMod))
}

func (t *Template) GetDirMode() os.FileMode {
	return os.FileMode(*t.DirMod.UnwrapOr(DefaultTemplateDirMod))
}

func (t *Template) GetIsFatal() bool {
	return *t.IsFatal.UnwrapOr(DefaultTemplateFatality)
}
