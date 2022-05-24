package testSampler

import (
	"testing"
)

type Sampler struct {
	t *testing.T
}

func New(t *testing.T) (sampler *Sampler) {
	sampler = &Sampler{
		t: t,
	}
	return
}
