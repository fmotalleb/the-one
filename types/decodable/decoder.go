package decodable

import (
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

func BuildDecoder(item any) (*mapstructure.Decoder, error) {
	hook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		DecodeHookFunc(),
	)

	decoderConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &item,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		DecodeHook:       hook,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	return decoder, err
}

func Decode[T any](input any, target *T) error {
	log := log().Named("Decode")
	decoder, err := BuildDecoder(target)
	if err != nil {
		log.Error("Decoder creation failed", zap.Error(err))
		return err
	}
	if err := decoder.Decode(input); err != nil {
		log.Error("Decode failed", zap.Error(err))
		return err
	}
	return nil
}
