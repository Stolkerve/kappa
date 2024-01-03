package types

import "time"

type RuntimeResult struct {
	ResponseWrapper ResponseWrapper
	Stdout          string
	Stderr          string
	Memory          uint32
	Duration        time.Duration
}
