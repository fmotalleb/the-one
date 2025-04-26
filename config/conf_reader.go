package config

import (
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/types/decodable"
)

// Decode map into Config struct.
func DecodeConfig(input map[string]any) (Config, error) {
	var cfg Config
	log().Info("Decoding config from map")
	if err := decodable.Decode(input, &cfg); err != nil {
		log().Error("Decode failed", zap.Error(err))
		return cfg, err
	}
	return cfg, nil
}
