package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/types/option"
)

type Config struct {
	Includes []string           `mapstructure:"include,omitempty"`
	Services map[string]Service `mapstructure:"services,omitempty"`
}

var logger = logging.GetLogger("config")

func DeepMerge(dst, src map[string]any) (map[string]any, error) {
	err := mergo.Merge(&dst, src, mergo.WithAppendSlice)
	return dst, err
}

// Load and merge config into map[string]any
func ReadAndMergeConfig(path string) (map[string]any, error) {
	logger.Info("Reading config", zap.String("path", path))
	visited := map[string]bool{}
	return readRecursive(path, visited)
}

func readRecursive(path string, visited map[string]bool) (map[string]any, error) {
	if visited[path] {
		logger.Error("Circular include detected", zap.String("path", path))
		return nil, fmt.Errorf("circular include detected: %s", path)
	}
	visited[path] = true

	v := viper.New()
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	v.SetConfigType(ext)
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		logger.Error("Failed to read config", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	raw := make(map[string]any)
	if err := v.Unmarshal(&raw); err != nil {
		logger.Error("Unmarshal failed", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	includes, _ := raw["include"].([]any)
	for _, inc := range includes {
		if incPath, ok := inc.(string); ok {
			logger.Info("Processing include", zap.String("include", incPath))
			childRaw, err := readRecursive(incPath, visited)
			if err != nil {
				return nil, err
			}
			raw, err = DeepMerge(childRaw, raw)
			if err != nil {
				logger.Error("Deep merge failed", zap.String("path", incPath), zap.Error(err))
				return nil, err
			}
		}
	}

	return raw, nil
}

// Decode map into Config struct
func DecodeConfig(input map[string]any) (Config, error) {
	var cfg Config
	logger.Info("Decoding config from map")
	hook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		option.DecodeHookFunc(),
	)

	decoderConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &cfg,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook:       hook,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		logger.Error("Decoder creation failed", zap.Error(err))
		return cfg, err
	}
	if err := decoder.Decode(input); err != nil {
		logger.Error("Decode failed", zap.Error(err))
		return cfg, err
	}
	return cfg, nil
}
