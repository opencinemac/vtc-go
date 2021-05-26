package testdata_test

import (
	"github.com/opencinemac/vtc-go/pkg/internal/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManyBasicEditsData(t *testing.T) {
	assert.Equal(t, int64(24), testdata.ManyBasicEditsData.StartTime.Timebase)
}
