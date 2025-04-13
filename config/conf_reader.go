package config

import (
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/types/option"
)

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
