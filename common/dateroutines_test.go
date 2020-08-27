package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcLongEpochFromEndDate(t *testing.T) {

	expectedResult := int64(1609477199)
	result := CalcLongEpochFromEndDate(2020, 12, 31)
	assert.Equal(t, expectedResult, result)

}
