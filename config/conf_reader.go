package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"dario.cat/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/types/option"
)

var log = sync.OnceValue(
	func() *zap.Logger {
		return logging.GetLogger("core.config")
	},
)

func DeepMerge(dst, src map[string]any) (map[string]any, error) {
	err := mergo.Merge(&dst, src, mergo.WithAppendSlice)
	return dst, err
}

func ReadAndMergeConfig(pattern string) (map[string]any, error) {
	log().Info("Starting config load", zap.String("pattern", pattern))
	return mergeFromPattern(pattern, map[string]bool{})
}

func mergeFromPattern(pattern string, visited map[string]bool) (map[string]any, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log().Error("Invalid glob pattern", zap.String("pattern", pattern), zap.Error(err))
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}
	if len(files) == 0 {
		log().Warn("No config files matched pattern", zap.String("pattern", pattern))
	}

	result := make(map[string]any)
	for _, file := range files {
		log().Debug("Merging config file", zap.String("file", file))
		conf, err := readAndResolveIncludes(file, visited)
		if err != nil {
			log().Error("Failed to read and merge includes", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		result, err = DeepMerge(result, conf)
		if err != nil {
			log().Error("Deep merge failed", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		log().Debug("Merged successfully", zap.String("file", file))
	}
	return result, nil
}

func readAndResolveIncludes(path string, visited map[string]bool) (map[string]any, error) {
	if visited[path] {
		log().Warn("Circular include detected", zap.String("path", path))
		return make(map[string]any), nil
	}
	log().Info("Reading config file", zap.String("path", path))
	visited[path] = true

	v := viper.New()
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(path), "."))
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		log().Error("Failed to read config", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	raw := make(map[string]any)
	if err := v.Unmarshal(&raw); err != nil {
		log().Error("Failed to unmarshal config", zap.String("path", path), zap.Error(err))
		return nil, err
	}
	log().Debug("Parsed config", zap.String("path", path))

	if includes, ok := raw["include"].([]any); ok {
		for _, inc := range includes {
			if pattern, ok := inc.(string); ok {
				log().Info("Processing include", zap.String("from", path), zap.String("pattern", pattern))
				included, err := mergeFromPattern(pattern, visited)
				if err != nil {
					log().Error("Failed to process include", zap.String("pattern", pattern), zap.Error(err))
					return nil, err
				}
				raw, err = DeepMerge(included, raw)
				if err != nil {
					log().Error("Deep merge failed during include", zap.String("pattern", pattern), zap.Error(err))
					return nil, err
				}
				log().Debug("Include merged", zap.String("pattern", pattern))
			}
		}
	}

	return raw, nil
}

// Decode map into Config struct.
func DecodeConfig(input map[string]any) (Config, error) {
	var cfg Config
	log().Info("Decoding config from map")
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
		log().Error("Decoder creation failed", zap.Error(err))
		return cfg, err
	}
	if err := decoder.Decode(input); err != nil {
		log().Error("Decode failed", zap.Error(err))
		return cfg, err
	}
	return cfg, nil
}
