package domain

import "time"

type Nower interface {
	Now() time.Time
}
