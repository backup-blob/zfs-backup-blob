package main_test

import (
	"context"
	"errors"
	"github.com/cucumber/godog"
	"io"
	"os"
)

const ConfigKey = "configKey"

func loadConfig(ctx context.Context, fileName string) (context.Context, error) {
	f, err := os.Open("./fixtures/" + fileName)
	if err != nil {
		return ctx, err
	}
	data, errRA := io.ReadAll(f)
	if errRA != nil {
		return ctx, errRA
	}
	return context.WithValue(ctx, ConfigKey, data), nil
}

func replaceConfigPlaceholder(ctx context.Context) (context.Context, error) {
	cfg, ok := ctx.Value(ConfigKey).([]byte)
	if ok {
		cfgStr := string(cfg)
		cfgReplaced, err := replacePlaceholder(ctx, cfgStr)
		if err != nil {
			return ctx, err
		}
		return context.WithValue(ctx, ConfigKey, []byte(cfgReplaced)), nil
	}
	return ctx, errors.New("did not replace")
}

func persistConfig(ctx context.Context, name string) (context.Context, error) {
	cfg, ok := ctx.Value(ConfigKey).([]byte)
	if ok {
		dir, err := os.MkdirTemp("", "example")
		if err != nil {
			return ctx, err
		}
		filePath := dir + "/" + "config.yaml"
		err = os.WriteFile(filePath, cfg, 0777)
		if err != nil {
			return ctx, err
		}
		return setPlaceholder(ctx, name, filePath)
	}
	return ctx, errors.New("persist config failed")
}

func InitConfig(sc *godog.ScenarioContext) {
	sc.Step(`the config is loaded from \'([_a-zA-Z-\/0-9@\.]+)\'`, loadConfig)
	sc.Step(`the config is persisted at \'([<>_a-zA-Z-\/0-9@\.]+)\'`, persistConfig)
	sc.Step(`the placeholders are replaced in config`, replaceConfigPlaceholder)
}
