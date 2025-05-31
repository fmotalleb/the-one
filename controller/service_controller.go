package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/renderer"
)

func Boot(ctx context.Context, cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", data)

	reshape, err := cfg.GetServices()
	if err != nil {
		return errors.Join(
			EngineBootError,
			err,
		)
	}
	data, err = json.Marshal(reshape)
	if err != nil {
		return errors.Join(
			EngineBootError,
			err,
		)
	}
	fmt.Printf("%s\n", data)

	// Compile Templates
	for _, t := range cfg.Templates {
		if err := renderer.RenderTemplates(&t); err != nil && t.GetIsFatal() {
			return errors.Join(
				EngineBootError,
				err,
			)
		}
	}

	return nil
}
