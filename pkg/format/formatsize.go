package format

import "github.com/docker/go-units"

func Size(size *int64) string {
	if size == nil {
		return "N/A"
	}

	return units.HumanSize(float64(*size))
}
