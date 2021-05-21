package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFakeRange(t *testing.T) {
	type Borne struct {
		min int64
		max int64
	}
	testingList := []Borne{
		{
			min: 0,
			max: 1,
		},
		{
			min: 1,
			max: 5,
		},
		{
			min: 0,
			max: 5000,
		},
		{
			min: 100,
			max: 101,
		},
		{
			min: 10,
			max: 11,
		},
	}
	for _, elem := range testingList {
		fakeRange := FakeRange(elem.min, elem.max)
		require.True(t, len(fakeRange) >= int(elem.min) && len(fakeRange) <= int(elem.max))
	}
}
