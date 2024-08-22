package main_test

import (
	"context"
	"errors"
	"strings"
)

const PlaceholderKey = "placeholderKey"

type Placeholder = map[string]string

func setPlaceholder(ctx context.Context, key, value string) (context.Context, error) {
	v, ok := ctx.Value(PlaceholderKey).(Placeholder)
	if !ok {
		v = map[string]string{}
	}
	v[key] = value
	return context.WithValue(ctx, PlaceholderKey, v), nil
}

func replacePlaceholder(ctx context.Context, str string) (string, error) {
	placeholder, okP := ctx.Value(PlaceholderKey).(Placeholder)
	if okP {
		for key, val := range placeholder {
			str = strings.ReplaceAll(str, key, val)
		}
		return str, nil
	}
	return "", errors.New("failed to replace")
}
