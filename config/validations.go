package config

import (
	"os"

	"github.com/go-playground/validator/v10"
)

// Custom validator function.
func workingDirValidator(fl validator.FieldLevel) bool {
	dir := fl.Field().String()
	if dir == "" {
		// Empty string is allowed
		return true
	}
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
