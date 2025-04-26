package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/fmotalleb/the-one/config"
)

func Boot(ctx context.Context, cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", data)

	reshape, err := cfg.GetServices()
	if err != nil {
		return err
	}
	data, err = json.Marshal(reshape)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)
	return nil
}
